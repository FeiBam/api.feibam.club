// 定义语言枚举
enum article_lang_t {
    ZH = 'zh', // 中文
    EN = 'en', // 英文
    JP = 'jp', // 日文
}

// 定义驱动类型枚举
const enum article_model_driver_type_t {
    MD,   // Markdown 驱动
    SQL,  // SQL 驱动
    JSON  // JSON 驱动
}

interface article_dirier_md_config_t {
    file_path : string
}

interface article_driver_sql_config_t {
    db_path: string;
    db_mode? : "read" | "write" | "create" | undefined
}

// 定义文章模型接口
interface article_model_t {
    id?: number;                // 唯一标识符
    title: string;             // 文章标题
    introduction: string;      // 文章简介
    tags: string[];            // 标签
    subject: string;           // 主题内容
    create_at: Date;           // 创建时间
    lang: article_lang_t;      // 语言
    links: URL[];              // 相关链接
}

// 驱动通用功能接口
interface article_driver_t {
    format(article: article_model_t): string | Array<unknown>; // 格式化数据为特定格式
    save(article: article_model_t): Promise<{ success: boolean; error?: Error }>; // 保存文件
    load(input: string | number): Promise<{ success: article_model_t | null; error?: Error }>; // 加载文件
    deformat(binary: Uint8Array | string | unknown): article_model_t; // 解析特定格式数据
}

// 配置类型接口
interface article_model_config_t {
    driver_type: article_model_driver_type_t; // 驱动类型
}

export { 
    article_lang_t,
    article_model_driver_type_t,
}

// 导出
export type { 
    article_model_t,
    article_driver_t,
    article_model_config_t,
    article_driver_sql_config_t,
    article_dirier_md_config_t
};
