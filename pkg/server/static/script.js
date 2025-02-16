function getSelectedValues(className) {
  return Array.from(document.getElementsByClassName(className))
    .filter((e) => e.checked)
    .map((e) => e.value);
}

function feedUrlIsValid() {
  const selectedRegions = getSelectedValues("region-select");
  const selectedInstanceTypes = getSelectedValues("instance-type-select");
  const selectedOperatingSystems = getSelectedValues("operating-system-select");

  return (
    selectedRegions.length > 0 &&
    selectedInstanceTypes.length > 0 &&
    selectedOperatingSystems.length > 0
  );
}

function generateFeedUrl() {
  const selectedRegions = getSelectedValues("region-select");
  const selectedInstanceTypes = getSelectedValues("instance-type-select");
  const selectedOperatingSystems = getSelectedValues("operating-system-select");

  const feedUrl = new URL("/feed", window.location.origin);
  feedUrl.searchParams.append("regions", selectedRegions.join(","));
  feedUrl.searchParams.append(
    "instance_types",
    selectedInstanceTypes.join(",")
  );
  feedUrl.searchParams.append(
    "operating_systems",
    selectedOperatingSystems.join(",")
  );

  return feedUrl.href;
}

function updateFeedUrlDisplay() {
  const errorDisplay = document.getElementById("feed-url-error-display");
  const feedContainer = document.getElementById("feed-url-container");
  if (!feedUrlIsValid()) {
    errorDisplay.style.display = "flex";
    feedContainer.style.display = "none";
  } else {
    errorDisplay.style.display = "none";
    feedContainer.style.display = "flex";
  }

  const feedUrl = generateFeedUrl();
  document.getElementById("feed-url").innerText = feedUrl;
}

for (const e of document.getElementsByClassName("feed-option")) {
  e.addEventListener("change", updateFeedUrlDisplay);
}

document.getElementById("copy-feed-button").addEventListener("click", () => {
  navigator.clipboard.writeText(generateFeedUrl());
});

updateFeedUrlDisplay();
