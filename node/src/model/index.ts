
import { Sequelize } from "npm:sequelize"

const { config } = await import("@app")

console.log(config)

const sequelize = new Sequelize({
    dialect: 'sqlite',
    storage: config.sql_storage_path
});

  /*
ArticleModel.belongsToMany(LinkModel, { through:"ArticleLink" })
LinkModel.belongsToMany(ArticleModel, { through:"ArticleLink" })
  
ArticleModel.belongsToMany(TagModel, { through:"ArticleTag" })
TagModel.belongsToMany(ArticleModel, { through:"ArticleTag" }) 
  */
  
export {
    sequelize,
}