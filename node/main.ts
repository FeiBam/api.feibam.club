
import { config_t } from "./src/types/config.ts";
import { Logger } from "jsr:@deno-library/logger";
import { parse } from "https://deno.land/std@0.201.0/flags/mod.ts";
import { server } from "./src/index.ts";

const args = parse(Deno.args);
const decoder = new TextDecoder("utf-8");

const config : config_t  = JSON.parse(decoder.decode(Deno.readFileSync('./config.json')))
const logger = new Logger();

await logger.initFileLogger(config.log_storage_path, {
  rotate: true,
  maxBytes: 9999999999,
  maxBackupCount: 7
})

const { sequelize } = await import("./src/model/index.ts")

console.log(config)

if (args.h || args.help){
  helpCommandHandler()
}

if (args._[0] === "run"){
  runCommandHandler(args._[1] as string, args)
}

if (args._[0] === 'sync'){
  syncCommandHandler(args._[1] as string, args)
}


function helpCommandHandler(){
  return console.log(`
Usage:
  deno run main.js [command] [options]

Commands:
  run <subcommand>   Start a specific service
    server           Start the server

  sync <subcommand>  Synchronize data or configuration
    db               Synchronize the database
    md               Synchronize markdown to database

  -h, -help          Show help information

Options:
  --port <number>    Specify the port for the server (default: 3000)

Examples:
  deno run main.js run server --port 8080
  deno run main.js sync db
  deno run main.js -h
`)
}

function runCommandHandler(type: string, args: { [x: string]: any; _: Array<string | number>; }){
  let listen_port = config.port
  if(args.port){
    if (args.port > 65535 || args.port < 1) {
      console.log("Error: The port number must be between 1 and 65535.");
      return;
    }
    listen_port = args.port
  }
  if(type === "server"){
    server.listen({
      port : listen_port
    })
  }
}

async function syncCommandHandler(type : string ,args: { [x: string]: any; _: Array<string | number>; }){

}


export {
  config,
  logger
}