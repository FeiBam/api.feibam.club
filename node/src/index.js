const Koa = require('koa')
const fs = require('fs')
const KoaBody = require('koa-bodyparser')
const repl = require('repl');

const {  loggerMiddleware, logger } = require("./loggerMiddleware.js")
const { ArticleModel, ArticleMD ,ArticlesControl} = require("./unity/Article.js")

const { Router } = require('./routes/main.js')

const Port = process.env.NODE_ENV !== 'development' ? `8000` : '80';

const app  = new Koa()

let config = getConfig()

function getConfig(){
    let configJson = fs.readFileSync('./config.json')
    return JSON.parse(configJson)
}

app.use(loggerMiddleware)

app.use(async (ctx,next)=>{
    ctx.set('Access-Control-Allow-Origin','*')
    ctx.set('Access-Control-Allow-Method','GET,POST')
    ctx.set('Access-Control-Allow-Headers','*')
    if(ctx.request.method === 'OPTIONS'){
        return ctx.status = 204
    }
    await next()
})

app.use(async (ctx,next)=>{
    if(ctx.originalUrl === '/commond/reload'){
        config = getConfig()
        return ctx.body = 'success'
    }
    ctx.state.globalConfig = config
    await next()
})


app.use(async (ctx,next)=>{
    const articlesControl = new ArticlesControl(ctx.state.globalConfig.Article.ArticlePath, ArticleMD)
    ctx.state.articlesControl = articlesControl
    await next()
})

app.use(KoaBody())

app.use(Router.allowedMethods())
app.use(Router.routes())

app.listen(80)

