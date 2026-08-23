function toggleActiveTeam(button) {
  const tabPanel = button.closest("#tab-panel");
  tabPanel
    .querySelectorAll("button")
    .forEach((btn) => btn.classList.remove("tab-active"));
  button.classList.add("tab-active");
}
