return {
  {
    "neovim/nvim-lspconfig",
    opts = {
      servers = {
        ruff = {
          cmd_env = { RUFF_TRACE = "messages" },
          init_options = {
            settings = {
              logLevel = "error",
              -- Configurações adicionais do Ruff
              lineLength = 88,
              lint = {
                -- Regras específicas que você quer ativar
                extendSelect = { "I", "E", "F" },
              },
              format = {
                -- Configurações de formatação
                preview = false,
              },
            },
          },
          keys = {
            {
              "<leader>co",
              function()
                vim.lsp.buf.code_action({
                  apply = true,
                  context = {
                    only = { "source.organizeImports" },
                  },
                })
              end,
              desc = "Organize Imports",
            },
            {
              "<leader>cf",
              function()
                vim.lsp.buf.code_action({
                  apply = true,
                  context = {
                    only = { "source.fixAll.ruff" },
                  },
                })
              end,
              desc = "Fix All (Ruff)",
            },
          },
        },
      },
      setup = {
        ruff = function()
          -- Desabilita o hover do Ruff para priorizar outros LSPs (ty, Pyright, etc)
          require("lazyvim.util").lsp.on_attach(function(client, _)
            if client.name == "ruff" then
              client.server_capabilities.hoverProvider = false
            end
          end)
        end,
      },
    },
  },
}
