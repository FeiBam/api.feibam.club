import { Model, DataTypes } from "npm:sequelize";
import { sequelize } from "./index.ts";

class Tag extends Model {}

Tag.init({
    id: {
        type: DataTypes.INTEGER,
        autoIncrement: true,
        primaryKey: true,
    },
    name: {
        type: DataTypes.STRING,
        allowNull: false
    }
},{
    sequelize: sequelize,
    modelName: "Tag",
    paranoid: true,
    deletedAt: 'deletedAt',
})


export { 
    Tag
}