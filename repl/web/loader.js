// Set the default global log for use by wasm_exec.js
go_log = console.log;

var useWasm = location.href.includes("?wasm");

console.log("useWasm =", useWasm);

var script = document.createElement('script');
if (useWasm) {
    script.src = "wasm_exec.js";
    script.onload = function () {
         // polyfill
        if (!WebAssembly.instantiateStreaming) {
            WebAssembly.instantiateStreaming = async (resp, importObject) => {
	        const source = await (await resp).arrayBuffer();
	        return await WebAssembly.instantiate(source, importObject);
            };
        }
    
        const go = new Go();
        let mod, inst;
        WebAssembly.instantiateStreaming(fetch("gpython.wasm"), go.importObject).then((result) => {
            mod = result.module;
            inst = result.instance;
            run();
        });
    
        async function run() {
	    console.clear();
            await go.run(inst);
            inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
        }
    };
} else {
    script.src = "gpython.js";
}
document.head.appendChild(script);
