// WORKER_JS_CONTENT - Main worker logic with configurable paths
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
    const url = new URL(request.url);

    // Route API requests to Go WASM handlers - configurable prefix
    if (url.pathname.startsWith("API_ROUTE_PREFIX")) {
      const binding = {};
      await run(createRuntimeContext({ env, ctx, binding }));
      return binding.handleRequest(request);
    }

    // Serve static assets for all other requests
    // CRITICAL: Without this, static files will not be served
    return env.ASSETS.fetch(request);
  },
};