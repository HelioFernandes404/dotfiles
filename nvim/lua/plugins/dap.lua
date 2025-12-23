return {
  {
    "mfussenegger/nvim-dap",
    config = function()
      -- Carrega launch.json automaticamente ao iniciar debug
      local dap = require("dap")

      -- Hook que carrega launch.json antes de continuar
      local continue = dap.continue
      dap.continue = function()
        if vim.fn.filereadable('.vscode/launch.json') == 1 then
          require('dap.ext.vscode').load_launchjs(nil, { debugpy = { 'python' } })
        end
        continue()
      end
    end,
  },
}
