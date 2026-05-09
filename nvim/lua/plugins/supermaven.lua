-- Run :SupermavenUseFree after installation to use the free tier
return {
  {
    "supermaven-inc/supermaven-nvim",
    event = "InsertEnter",
    opts = {
      keymaps = {
        accept_total = "<C-y>",
        accept_partial = "<C-w>",
        ignore = "<C-e>",
      },
    },
  },
  {
    "Saghen/blink.cmp",
    opts = {
      ghost_text = false,
    },
  },
}