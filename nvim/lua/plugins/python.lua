return {
  {
    "neovim/nvim-lspconfig",
    opts = {
      servers = {
        -- Desabilitar pyright (redundante com ty)
        pyright = { enabled = false },
        basedpyright = { enabled = false },

        -- ty: type checker + LSP + IntelliSense
        ty = {
          settings = {
            ty = {
              completions = { autoImport = true },
              diagnosticMode = "workspace",
              inlayHints = {
                callArgumentNames = true,
                variableTypes = true,
              },
            },
          },
        },

        -- ruff: linting + formatting (complementar ao ty)
        ruff = {},
      },
    },
  },

  -- debugpy: apenas para debugging (DAP)
  {
    "mfussenegger/nvim-dap-python",
    dependencies = { "mfussenegger/nvim-dap" },
  },
}
