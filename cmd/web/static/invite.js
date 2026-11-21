// Invite link on the Users page. The server sends only the path of the link;
// the browser adds its own address so the link works on any host.

// A swapped-in invite arrives with an empty field: fill it in.
document.addEventListener("htmx:afterSwap", () => {
  const invite = document.querySelector("[data-invite-path]");
  if (!invite) return;
  const field = invite.querySelector("input");
  field.value = location.origin + invite.dataset.invitePath;
});

// Share opens the phone's share sheet. Where there is none, select the text
// so it can be copied.
document.addEventListener("click", (event) => {
  const button = event.target.closest("[data-invite-share]");
  if (!button) return;
  const field = button.parentElement.querySelector("input");
  if (navigator.share) {
    navigator.share({ url: field.value }).catch(() => {});
    return;
  }
  field.focus();
  field.select();
  field.setSelectionRange(0, field.value.length);
});
