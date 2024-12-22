import { Model, DataTypes } from "npm:sequelize";
import { sequelize } from "./index.ts";

class Link extends Model {}

Link.init({
    id: {
        type: DataTypes.INTEGER,
        autoIncrement: true,
        primaryKey: true,
    },
    name: { 
        type: DataTypes.STRING,
        allowNull: false
    },
    url: {
        type: DataTypes.STRING,
        allowNull: false
    }
},{
    sequelize: sequelize,
    modelName: "Link",
    paranoid: true,
    deletedAt: 'deletedAt',
})


export {
    Link
}