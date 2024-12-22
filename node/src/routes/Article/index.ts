
import { Router } from "jsr:@oak/oak/router";
import { respondHandel } from "../../unity/respondHandel.js"
import { ArticleModel } from "../../unity/Article.ts";
import { article_driver_sql_config_t, article_model_driver_type_t } from "../../types/articles-types.ts";


const router = new Router();
const driver_config : article_driver_sql_config_t = {
    db_path:'./test.db'
}
const driver = ArticleModel.getDriver(article_model_driver_type_t.SQL, driver_config)

router.use(async (ctx)=>{
    ctx.state.article_driver = driver
})

router.get('/:lang/:id',async (ctx) => {
    let article_id = ctx.params.id
    let article_lang = ctx.params.lang
})



export { 
    router
}