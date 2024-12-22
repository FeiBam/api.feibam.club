
import { Article as ArticleModel } from "./Article.ts"
import { Link as LinkModel } from "./Link.ts"
import { Tag as TagModel } from "./Tag.ts"

ArticleModel.belongsToMany(LinkModel, { through:"ArticleLink" })
LinkModel.belongsToMany(ArticleModel, { through:"ArticleLink" })

ArticleModel.belongsToMany(TagModel, { through:"ArticleTag" })
TagModel.belongsToMany(ArticleModel, { through:"ArticleTag" }) 


export {
    ArticleModel,
    LinkModel,
    TagModel
}