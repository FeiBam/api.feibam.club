import { Application } from "jsr:@oak/oak/application";
import { Router } from "jsr:@oak/oak/router";

import { config } from "@app";
import { logger } from "@app";


const router = new Router();


router.get("/:lang/:id", (ctx) => {
  ctx.response.body = ctx.request
  console.log(ctx.params)
  
});

router.get("/",(ctx)=>{
  ctx.response.body = "Hello World"
})



const server = new Application();

server.use((ctx,next)=>{
  ctx.state.global_config = config // 加入config 文件到 服务器ctx内
  ctx.state.logger = logger
  return next()
})


server.use(router.routes());
server.use(router.allowedMethods());


export {
  server
}

