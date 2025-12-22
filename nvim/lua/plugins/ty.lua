return {
  {
    "neovim/nvim-lspconfig",
    opts = {
      servers = {
        ty = {
          -- Comando para iniciar o servidor
          cmd = { "ty", "server" },

          -- Tipos de arquivo onde o ty vai rodar
          filetypes = { "python" },

          -- Configurações do servidor
          settings = {
            ty = {
              -- Diagnósticos em todo o workspace, não só arquivos abertos
              diagnosticMode = "workspace",

              -- Manter os language services habilitados (completion, hover, etc.)
              disableLanguageServices = false,

              -- Inlay hints (dicas inline de tipos)
              inlayHints = {
                variableTypes = true,     -- mostrar tipos de variáveis
                callArgumentNames = true, -- mostrar nomes de argumentos em chamadas
              },

              -- Completions (autocompletar)
              completions = {
                autoImport = true, -- sugerir imports automáticos
              },
            },
          },

          -- Opções de inicialização (opcional, para debug)
          init_options = {
            logLevel = "info", -- "trace" | "debug" | "info" | "warn" | "error"
            -- logFile = vim.fn.stdpath("cache") .. "/ty.log",  -- descomente para logar em arquivo
          },
        },
      },

      -- Setup customizado para o ty
      setup = {
        ty = function(_, opts)
          local lspconfig = require("lspconfig")
          local configs = require("lspconfig.configs")

          -- Registrar o ty como servidor customizado se ainda não existe
          if not configs.ty then
            configs.ty = {
              default_config = {
                cmd = opts.cmd,
                filetypes = opts.filetypes,
                root_dir = lspconfig.util.root_pattern(
                  "pyproject.toml",
                  "setup.py",
                  "setup.cfg",
                  "requirements.txt",
                  "Pipfile",
                  ".git"
                ),
                single_file_support = true,
                settings = opts.settings,
                init_options = opts.init_options,
              },
            }
          end

          lspconfig.ty.setup(opts)
        end,
      },
    },
  },
}
