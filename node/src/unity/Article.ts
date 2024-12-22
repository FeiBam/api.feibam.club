import {
    article_model_t,
    article_driver_t,
    article_model_driver_type_t,
    article_lang_t,
    article_driver_sql_config_t,
    article_dirier_md_config_t
} from '../types/articles-types.ts';

import matter from "npm:gray-matter";
import { DB, QueryParameterSet, SqliteError, SqliteOptions } from "https://deno.land/x/sqlite/mod.ts";

// Markdown 驱动实现
class MDDriver implements article_driver_t {
    private file_locationg : string | undefined
    constructor(driverConfig?: article_dirier_md_config_t){
        if(!driverConfig){
            this.file_locationg = undefined
        }
        else this.file_locationg = driverConfig.file_path
    }
    format(article: article_model_t): string {
        const frontMatter = matter.stringify('', {
            id: article.id,
            title: article.title,
            introduction: article.introduction,
            tags: article.tags,
            subject: article.subject,
            create_at: article.create_at.toISOString(),
            lang: article.lang,
            links: article.links.map(link => link.toString()),
        });
        return `${frontMatter}\n${article.subject}`;
    }

    async save(article: article_model_t, path? : string): Promise<{ success: boolean; error?: Error }> {
        try {
            if (!article.id){
                throw new Error("if you use MDDriver you must be set id")
            }
            if (!this.file_locationg || path){
                throw new Error("not have that path")
            }        
            const formattedText = this.format(article);
            await Deno.writeTextFile(this.file_locationg, formattedText);
            return { success: true };
        } catch (error) {
            return {
                success: false,
                error: error instanceof Error ? error : new Error("Unknown error occurred"),
            };
        }
    }

    async load(path: string): Promise<{ success: article_model_t | null; error?: Error }> {
        this.file_locationg = path

        try {
            const decoder = new TextDecoder("utf-8");
            const data = await Deno.readFile(path);
            const markdown = decoder.decode(data);
            const articleData = this.deformat(markdown);
            return { success: articleData };
        } catch (error) {
            return {
                success: null,
                error: error instanceof Error ? error : new Error("Unknown error occurred"),
            };
        }
    }

    deformat(binary: Uint8Array | string): article_model_t {
        const parsed = matter(binary);
        const metadata = parsed.data as Partial<article_model_t>;
        const content = parsed.content;

        if ((!metadata.id && metadata.id !== 0) || !metadata.title || !metadata.create_at || !metadata.lang) {
            throw new Error("Invalid article data format");
        }

        return {
            id: metadata.id,
            title: metadata.title,
            introduction: metadata.introduction || '',
            tags: metadata.tags || [],
            subject: content,
            create_at: new Date(metadata.create_at),
            lang: metadata.lang as article_lang_t,
            links: (metadata.links || []).map(link => new URL(link)),
        };
    }
}

class SQLDriver implements article_driver_t{
    
    private config : article_driver_sql_config_t
    private db_path : string
    private db : DB
    constructor (driverConfig : article_driver_sql_config_t){
        this.config = driverConfig
        this.db_path = this.config.db_path
        const db = new DB(this.db_path , { mode : this.config.db_mode });
        this.db = db
        db.execute(`
            CREATE TABLE IF NOT EXISTS people (
              id INTEGER PRIMARY KEY AUTOINCREMENT,
              title TEXT,
              introduction TEXT,
              tags BLOB,
              subject TEXT,
              create_at TEXT,
              lang INTEGER,
              links  BLOB
            )
          `);
    }
    deformat(binary: Uint8Array | string | unknown): article_model_t {
        const row = binary as unknown[];
        if (!Array.isArray(row)) {
            throw new Error("Invalid row format: expected an array");
        }
        return {
            id: row[0] as number, // ID
            title: row[1] as string, // 标题
            introduction: row[2] as string, // 简介
            tags: JSON.parse(row[3] as string), // Tags (JSON 解析)
            subject: row[4] as string, // 主题内容
            create_at: new Date(row[5] as string), // 创建日期
            lang: row[6] as article_lang_t, // 语言
            links: JSON.parse(row[7] as string), // Links (JSON 解析)
        };
    }
    format(article: article_model_t): Array<unknown> {
        return [
            article.title, // 标题
            article.introduction, // 简介
            JSON.stringify(article.tags), // Tags 转成 JSON 字符串
            article.subject, // 主题
            article.create_at.toISOString(), // 日期格式化为 ISO 字符串
            article.lang, // 语言
            JSON.stringify(article.links), // Links 转成 JSON 字符串
        ];
    }
    async load(id: string | number): Promise<{ success: article_model_t | null; error?: Error; }> {
      try {
        if (!(typeof id === "number" )) {
            const try_to_number = Number(id);
            if(Number.isNaN(try_to_number) || try_to_number <= 0 || !Number.isInteger(try_to_number)){
                return {
                    success: null,
                    error : new Error("parament type error")
                }
            }
            id = try_to_number
        }
        const result = this.db.query(`SELECT * FROM Article WHERE id = ?`, [id]);
        const formattedResult = this.deformat(result[0] as unknown[])
        return {
            success : formattedResult,
        }
      }catch (error){
        return {
            success: null,
            error: error instanceof Error ? error : new Error("Unknown error occurred"),
        };
      }
    }
    async save(article: article_model_t): Promise<{ success: boolean; error?: Error }> {
        try {
            const formatModel = this.format(article);
    
            this.db.query(
                `INSERT INTO Article (
                    title,
                    introduction,
                    tags,
                    subject,
                    create_at,
                    lang,
                    links
                ) VALUES (?, ?, ?, ?, ?, ?, ?);`,
                formatModel  as QueryParameterSet// 参数化查询
            );
    
            return { success: true };
        } catch (error) {
            return {
                success: false,
                error: error instanceof Error ? error : new Error("Unknown error occurred"),
            };
        }
    }
}

class SequlizeDriver implements article_driver_t {
}
// ArticleModel 主类
class ArticleModel implements article_model_t {
    id?: number;
    title: string;
    introduction: string;
    tags: string[];
    subject: string;
    create_at: Date;
    lang: article_lang_t;
    links: URL[];
    private driver: article_driver_t;

    constructor(article_model: article_model_t , driver: article_driver_t) {
        this.id = article_model.id;
        this.title = article_model.title;
        this.introduction = article_model.introduction;
        this.tags = article_model.tags;
        this.subject = article_model.subject;
        this.create_at = article_model.create_at;
        this.lang = article_model.lang;
        this.links = article_model.links;
        this.driver = driver;
    }

    format(): string | unknown {
        return this.driver.format(this);
    }

    async save(): Promise<{ success: boolean; error?: Error }> {
        return this.driver.save(this);
    }

    static async load(
        input: string | number,
        driverType: article_model_driver_type_t,
        driverConfig?: article_driver_sql_config_t | article_dirier_md_config_t 
    ): Promise<{ success: ArticleModel | null; error?: Error }> {
        const driver = ArticleModel.getDriver(driverType, driverConfig);
        const result = await driver.load(input);
        if (result.success) {
            return {
                success: new ArticleModel(result.success, driver),
            };
        }
        return {
            success: null,
            error: result.error,
        };
    }

    static getDriver(driverType: article_model_driver_type_t , driverConfig? : article_dirier_md_config_t | article_driver_sql_config_t): article_driver_t {
        switch (driverType) {
            case article_model_driver_type_t.MD:
                return new MDDriver(driverConfig as article_dirier_md_config_t);
            case article_model_driver_type_t.SQL:
                return new SQLDriver(driverConfig as article_driver_sql_config_t)
            case article_model_driver_type_t.JSON:
                throw new Error("JSON driver not implemented yet.");
            default:
                throw new Error("Invalid driver type.");
        }
    }
}

export {
    ArticleModel
}