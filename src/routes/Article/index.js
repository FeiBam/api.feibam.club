const KoaRouter = require('@koa/router') //导入Koa-Router 库
const Router = new KoaRouter()

const { respondHandel } = require('../../unity/respondHandel')

Router.post('/',async (ctx)=>{
    let RequestData = ctx.request.body
    let ResponseObject = {}
    let article = await ctx.state.articlesControl.getArticleByLangAndId(RequestData.id,RequestData.language)
    ResponseObject['id'] = Number(article.id)
    ResponseObject['Account'] = {
        Name:'Fei_Bam'
    }
    ResponseObject['Article'] = article
    return respondHandel.success(ctx,ResponseObject,'ok')
})

module.exports = Router