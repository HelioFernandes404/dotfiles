return {
  -- adicionar tokyonight
  {
    "folke/tokyonight.nvim",
    lazy = false,
    priority = 1000,
    opts = {
      style = "night", -- Opções: "storm", "moon", "night", "day"
    },
  },

  -- Configurar LazyVim para carregar o tokyonight
  {
    "LazyVim/LazyVim",
    opts = {
      colorscheme = "tokyonight",
    },
  },
}
