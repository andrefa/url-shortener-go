document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("shorten-form");
  const shortenedLink = document.getElementById("shortened-link");
  const shortenedUrl = document.getElementById("shortened-url");

  form.addEventListener("submit", async (event) => {
    event.preventDefault();

    const urlInput = document.getElementById("url");
    const url = urlInput.value;

    try {
      const response = await fetch("/api/v1/shorten", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url }),
      });

      if (!response.ok) {
        throw new Error("Failed to shorten URL");
      }

      const data = await response.json();
      shortenedLink.classList.remove("hidden");
      shortenedUrl.textContent = data.shortenedUrl;
      shortenedUrl.href = data.shortenedUrl;
    } catch (error) {
      alert("Error: " + error.message);
    }
  });

  document.getElementById("copy-btn").addEventListener("click", () => {
    const shortenedUrl = document.getElementById("shortened-url").href;
    navigator.clipboard.writeText(shortenedUrl).then(() => {
      alert("Shortened URL copied to clipboard!");
    });
  });
});
