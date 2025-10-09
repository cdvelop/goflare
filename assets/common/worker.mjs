import "./wasm_exec.js";
import { createRuntimeContext, loadModule } from "./runtime.mjs";

let mod;

globalThis.tryCatch = (fn) => {
  try {
    return {
      result: fn(),
    };
  } catch (e) {
    return {
      error: e,
    };
  }
};

async function run(ctx) {
  if (mod === undefined) {
    mod = await loadModule();
  }
  const go = new Go();

  let ready;
  const readyPromise = new Promise((resolve) => {
    ready = resolve;
  });
  const instance = new WebAssembly.Instance(mod, {
    ...go.importObject,
    workers: {
      ready: () => {
        ready();
      },
    },
  });
  go.run(instance, ctx);
  await readyPromise;
}

export default {
  async fetch(request, env, ctx) {
    const binding = {};
    await run(createRuntimeContext({ env, ctx, binding }));
    return binding.handleRequest(request);
  },

  async scheduled(event, env, ctx) {
    const binding = {};
    await run(createRuntimeContext({ env, ctx, binding }));
    return binding.handleScheduled(event);
  },

  async queue(batch, env, ctx) {
    const binding = {};
    await run(createRuntimeContext({ env, ctx, binding }));
    return binding.handleQueue(batch);
  },
};