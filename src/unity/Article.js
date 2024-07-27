const fs = require('fs');
const matter = require('gray-matter');
const MarkdownIt = require('markdown-it');


class ArticleModel {
    constructor(title = '', Introduction = '', Tags = [], Subject = '') {
        if (new.target === ArticleModel) {
            throw new TypeError("Cannot construct Abstract instances directly");
        }
        this.Title = title;
        this.Introduction = Introduction;
        this.Tags = Tags;
        this.Subject = Subject;
        this.create_At = new Date().getTime();
        this.id = null;
        this.lang = null;
        this.link = {};
    }

    format() {
        throw new Error("Method 'format()' must be implemented.");
    }

    deformat() {
        throw new Error("Method 'deformat()' must be implemented.");
    }

    save(path) {
        throw new Error("Method 'save()' must be implemented.");
    }

    static load(path) {
        throw new Error("Method 'load()' must be implemented.");
    }

    set_title(title) {
        this.Title = title;
    }

    set_Introduction(Introduction) {
        this.Introduction = Introduction;
    }

    set_Tags(Tags) {
        this.Tags = Tags;
    }

    set_Subject(Subject) {
        this.Subject = Subject;
    }

    set_create_At(Time) {
        this.create_At = Time;
    }

    add_Tag(Name) {
        this.Tags.push(Name);
    }

    remove_Tag(Name) {
        const index = this.Tags.indexOf(Name);
        if (index === -1) {
            return 'Not Have this Tag';
        }
        this.Tags.splice(index, 1);
    }

    get_CreateAt() {
        return this.create_At;
    }

    get_ArticleCreateDate() {
        const date = new Date(this.get_CreateAt());
        const year = date.getFullYear();
        const month = (date.getMonth() + 1).toString().padStart(2, '0');
        const day = date.getDate().toString().padStart(2, '0');
        return `${year}-${month}-${day}`; // 格式化日期为 YYYY-MM-DD
    }

    getLang() {
        return this.lang;
    }

    setLang(lang) {
        this.lang = lang;
    }

    getID() {
        return this.id;
    }

    setID(id) {
        this.id = id;
    }

    setOtherLangArticle(lang, ArticleID) {
        if (this.link.hasOwnProperty(lang)) throw Error('Already have this lang Article');
        this.link[lang] = ArticleID;
    }

    getOtherLangArticle(lang) {
        if (!this.link.hasOwnProperty(lang)) return false;
        return this.link[lang];
    }
}


class ArticleMD extends ArticleModel {
    format() {
        const frontMatter = matter.stringify('', {
            title: this.Title,
            introduction: this.Introduction,
            tags: this.Tags,
            subject: this.Subject,
            create_At: this.create_At,
            id: this.id,
            lang: this.lang,
            link: this.link
        });
        return `${frontMatter}\n${this.Subject}`;
    }

    deformat(markdown) {
        const parsed = matter(markdown);
        const data = parsed.data;
        const content = parsed.content;
        this.set_title(data.title);
        this.set_Introduction(data.introduction);
        this.set_Tags(data.tags);
        this.set_Subject(content.trim());
        this.set_create_At(data.create_At);
        this.setID(data.id);
        this.setLang(data.lang);
        this.link = data.link;
    }

    save(path) {
        const formattedText = this.format();
        fs.writeFileSync(path, formattedText);
    }

    static load(path) {
        if(!path.endsWith('.md')){
            path = `${path}.md`
        }
        const markdown = fs.readFileSync(path, 'utf8');
        const article = new ArticleMD();
        article.deformat(markdown);
        return article;
    }
}


class ArticleJSON extends ArticleModel {
    format() {
        return JSON.stringify(this);
    }

    deformat(json) {
        const data = JSON.parse(json);
        this.set_title(data.Title);
        this.set_Introduction(data.Introduction);
        this.set_Tags(data.Tags);
        this.set_Subject(data.Subject);
        this.set_create_At(data.create_At);
        this.setID(data.id);
        this.setLang(data.lang);
        this.link = data.link;
    }

    save(path) {
        const json = this.format();
        fs.writeFileSync(path, json);
    }

    static load(path) {
        const json = fs.readFileSync(path, 'utf8');
        const article = new ArticleJSON();
        console.log(json)
        article.deformat(json);
        return article;
    }
}


class ArticlesControl {
    #Articles_Path;
    #private_articles_tree;
    #Article_Class;

    constructor(path, ArticleClass) {
        this.#Article_Class = ArticleClass;
        this.#Articles_Path = path;
        this.#private_articles_tree = {};
        this.Need_Update = false;
        if (!fs.existsSync(path)) {
            fs.mkdirSync(path);
        }
    }

    getArticles_Tree() {
        return this.#private_articles_tree;
    }

    #onLang(lang) {
        this.#private_articles_tree[lang] = new Map()
    }

    async #readArticle(lang, ArticleFileName) {
        const ArticlePath = `${this.#Articles_Path}/${lang}/${ArticleFileName}`;
        return this.#Article_Class.load(ArticlePath);
    }

    async #Update_Tree(Lang, Article) {
        const ArticleCreateDate = Article.get_ArticleCreateDate();
        Article.createDate = ArticleCreateDate;
        this.#private_articles_tree[Lang].set(Article.getID() , Article)
        return this;
    }

    async save() {
        for (let lang of Object.keys(this.#private_articles_tree)) {
            for (let ArticleId of Object.keys(this.#private_articles_tree[lang]).filter((item) => item !== 'index')) {
                const path = `${this.#Articles_Path}/${lang}/${ArticleId}`;
                this.#private_articles_tree[lang][ArticleId].save(path);
            }
        }
        return true;
    }

    async getAll() {
        const dirList = fs.readdirSync(this.#Articles_Path);
        for (let Lang of dirList) {
            this.#onLang(Lang);
            const ArticleList = fs.readdirSync(`${this.#Articles_Path}/${Lang}`);
            for (let ArticleFileName of ArticleList) {
                const Article = await this.#readArticle(Lang, ArticleFileName);
                await this.#Update_Tree(Lang, Article);
            }
        }
        return this;
    }

    async getArticleByLangAndId(ArticleId, Lang) {
        const Article = await this.#readArticle(Lang, ArticleId);
        this.#onLang(Lang);
        await this.#Update_Tree(Lang, Article);
        return this.#private_articles_tree[Lang].get(ArticleId);
    }

    async getArticlesByLang(Lang) {
        this.#onLang(Lang);
        console.log(this.#private_articles_tree)
        const dirList = fs.readdirSync(`${this.#Articles_Path}/${Lang}`);
        console.log(dirList)
        for (let ArticleFileName of dirList) {
            const Article = await this.#readArticle(Lang, ArticleFileName);
            await this.#Update_Tree(Lang, Article);
        }
        return this.#private_articles_tree[Lang];
    }

    async createArticle(lang, ...args) {
        if (!lang) {
            throw new Error('You need Set Article Language');
        }
        this.Need_Update = true;
        await this.getAll();
        if (!this.hasLang(lang)) {
            this.newLang(lang);
        }
        const ArticleID = this.#private_articles_tree[lang].size;
        if (args.size === 1 && args[0] instanceof this.#Article_Class) {
            await this.#Update_Tree(lang, ArticleID, args[0]);
            return args[0]
        } else {
            const article = new this.#Article_Class(...args);
            await this.#Update_Tree(lang, ArticleID, Article);
            return article
        }
    }

    async getAllArticleNum() {
        await this.getAll();
        let ArticleNum = 0;
        for (let lang of Object.keys(this.#private_articles_tree)) {
            ArticleNum += this.#private_articles_tree[lang].size;
        }
        return ArticleNum;
    }

    async getArticleNumByLang(lang) {
        await this.getAll();
        if (!this.#private_articles_tree[lang]) return false;
        return this.#private_articles_tree[lang].size;
    }

    async getArticlesByTag(tag){
        await this.getAll();
        let filter = []

        for (const [lang, articlesMap] of Object.entries(this.#private_articles_tree)) {
            for (const article of articlesMap.values()) {
                if (article.Tags.includes(tag)) {
                    filter.push(article);
                }
            }
        }
        return filter
    }

    async getArticlesByTagWithLang(tag,lang) {
        await this.getAll();
        let filter = [];

        if (this.#private_articles_tree.hasOwnProperty(lang)) {
            const articlesMap = this.#private_articles_tree[lang];
            for (const article of articlesMap.values()) {
                if (article.Tags.includes(tag)) {
                    filter.push(article);
                }
            }
        }
        return filter;
    }


    hasLang(lang) {
        return Object.keys(this.#private_articles_tree).includes(lang);
    }

    async newLang(lang) {
        fs.mkdirSync(`${this.#Articles_Path}/${lang}`);
        this.#onLang(lang);
        return true;
    }

    LinkOtherLangArticle(selfArticle, targetArticle) {
        selfArticle.setOtherLangArticle(targetArticle.lang, targetArticle.id);
        targetArticle.setOtherLangArticle(selfArticle.lang, selfArticle.id);
    }

    async getOtherLangArticle(lang, Article) {
        const otherLangArticle = await this.getArticleByLangAndId(Article.getOtherLangArticle(lang), lang);
        return otherLangArticle;
    }

    static Mixin(FuncName, Call, targetClass) {
        if (targetClass[FuncName]) throw Error('the class already have this function');
        targetClass[FuncName] = Call;
    }
}

module.exports = {
    ArticleModel,
    ArticleMD,
    ArticleJSON,
    ArticlesControl
};
