import "./css/styles.css";

import "./components/navigation/navigation";
import "./components/content/content";

class PlatformApp extends HTMLElement {
  static observedAttributes = [];

  constructor() {
    super();
  }

  attributeChangedCallback(name, oldValue, newValue) {
    console.log(`Attribute ${name} has changed.`);
  }
  connectedCallback() {
    console.log("Platform connected");
  }
  disconnectedCallback() { }
  connectedMoveCallback() {
    console.log("Custom move-handling logic here.");
  }
}

if (!customElements.get("plc-platform")) {
  customElements.define("plc-platform", PlatformApp);
}

