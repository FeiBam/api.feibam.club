const KoaRouter = require('@koa/router') //导入Koa-Router 库
const Router = new KoaRouter()
const { ArticleModel, ArticlesControl} = require('../../unity/Article')
const { respondHandel } = require('../../unity/respondHandel')

Router.post('/getPageInfo',async (ctx)=>{
    let RequestData = ctx.request.body
    const articlesControl = ctx.state.articlesControl
    let ArticleNum = await articlesControl.getAllArticleNum()
    let ArticleNumofLang = await articlesControl.getArticleNumByLang(RequestData.language)
    if(!ArticleNumofLang) ArticleNumofLang = 0
    const ArticleLimit = ctx.state.globalConfig.Page.ArticleLimit
    return ctx.body = {
        PageArticleLimit:ArticleLimit,
        ArticleNum:ArticleNum,
        ArticleNumofLang
    }
})

Router.get('/page','/:lang/:index',async (ctx)=>{
    let RequestData = ctx.request.body
    const articlesControl = ctx.state.articlesControl
    const articles = await articlesControl.getArticlesByLang(ctx.params.lang)
    if(!articles){
        return ctx.body = {
            code:404,
            data:{},
            msg:'Not Found'
        }
    }
    return respondHandel.success(ctx,Array.from(articles.values()),'ok')
})

module.exports = Router