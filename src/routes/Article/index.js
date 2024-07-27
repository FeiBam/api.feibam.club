const KoaRouter = require('@koa/router') //导入Koa-Router 库
const Router = new KoaRouter()

const { respondHandel } = require('../../unity/respondHandel')

Router.post('/',async (ctx)=>{
    let requestData = ctx.request.body
    let responseObject = {}
    let article = await ctx.state.articlesControl.getArticleByLangAndId(requestData.id,requestData.language)
    responseObject['id'] = Number(article.id)
    responseObject['Account'] = {
        Name:'Fei_Bam'
    }
    responseObject['Article'] = article
    return respondHandel.success(ctx,responseObject,'ok')
})


Router.post(`/tag`,async (ctx)=>{
    let requestData = ctx.request.body
    let responseObject = {}
    let articles = await ctx.state.articlesControl.getArticlesByTagWithLang(requestData.tag,requestData.language);
    if(!articles){
        return ctx.body = {
            code:404,
            data:{},
            msg:'Not Found'
        }
    }
    return respondHandel.success(ctx,articles,'ok')

})

module.exports = Router