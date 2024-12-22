const fs = require('fs').promises; // 使用异步版本的fs
const path = require('path');
const log4js = require('log4js');
const Stream = require('stream');

const config = (function () {
    const configJson = require(path.resolve(__dirname, './config.json')); // 直接读取JSON文件
    return configJson;
})();

const Dates = new Date();
const NowMonth = String(Dates.getMonth() + 1).padStart(2, '0');
const NowDate = String(Dates.getDate()).padStart(2, '0');
const NowYear = Dates.getFullYear();

const logName = `${NowYear}-${NowMonth}-${NowDate}`;
const logDir = config.logPath;
let loggerDir = '';

(async function () {
    try {
        // 检查并创建日志目录
        if (!await fs.access(logDir).then(() => true).catch(() => false)) {
            await fs.mkdir(logDir, { recursive: true });
        }

        const files = await fs.readdir(logDir);

        loggerDir = path.join(logDir, `${logName}.log`);
        if (!files.includes(`${logName}.log`)) {
            await fs.writeFile(loggerDir, `This is Server Log Create at ${logName}\n`);
        }

        log4js.configure({
            appenders: {
                console: { type: 'console' },
                dateFile: { type: 'dateFile', filename: loggerDir, pattern: '-yyyy-MM-dd' }
            },
            categories: {
                default: {
                    appenders: ['console', 'dateFile'],
                    level: 'info'
                }
            }
        });
    } catch (error) {
        console.error('Error setting up logger:', error);
    }
})();

const logger = log4js.getLogger('[Default]');

const loggerMiddleware = async (ctx, next) => {
    const start = new Date();
    await next();
    const ms = new Date() - start;
    const remoteAddress = ctx.headers['x-forwarded-for'] || ctx.ip || ctx.ips ||
        (ctx.socket && (ctx.socket.remoteAddress || (ctx.socket.socket && ctx.socket.socket.remoteAddress)));
    let FileType = ctx.state.fileType;
    if (ctx.status === 404) {
        FileType = 'html';
    }
    let logText = `${ctx.method} ${ctx.status} ${decodeURI(ctx.url)} 响应数据: - ${remoteAddress} - ${ms}ms`;
    if (ctx.body instanceof Stream || ctx.body instanceof Buffer) {
        logText = `${ctx.method} ${ctx.status} ${decodeURI(ctx.url)} 文件类型 ${FileType} - ${remoteAddress} - ${ms}ms`;
        if (ctx.status === 206) {
            logText = `${ctx.method} ${ctx.status} ${decodeURI(ctx.url)} 请求数据范围： ${ctx.response.get('Content-Range')} 文件类型 ${FileType} - ${remoteAddress} - ${ms}ms`;
        }
    } else {
        logText = `${ctx.method} ${ctx.status} ${decodeURI(ctx.url)} - ${JSON.stringify(ctx.body)} - ${remoteAddress} - ${ms}ms`;
    }
    logger.info(logText);
};

module.exports = {
    loggerMiddleware,
    logger
};