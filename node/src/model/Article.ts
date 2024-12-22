import { Model, DataTypes } from "npm:sequelize";
import { sequelize } from "./index.ts";


class Article extends Model {}

Article.init({
    id: {
        type: DataTypes.INTEGER,
        autoIncrement: true,
        primaryKey: true,
    },
    title: { 
        type: DataTypes.STRING,
        allowNull: false
    },
    introduction: {
        type: DataTypes.STRING,
        allowNull: false
    },
    subject: {
        type: DataTypes.STRING,
        allowNull: false
    },
    lang: {
        type: DataTypes.INTEGER,
        allowNull: false
    }
},{
    sequelize: sequelize,
    modelName: "Article",
    paranoid: true,
    deletedAt: 'deletedAt',
})


export {
    Article
}