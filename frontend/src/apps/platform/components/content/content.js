class ContentComponent extends HTMLElement {
  static observedAttributes = [];

  constructor() {
    super();
  }

  attributeChangedCallback(name, oldValue, newValue) {
    console.log(`Attribute ${name} has changed.`);
  }
  connectedCallback() {
    console.log("Content connected");
  }
  disconnectedCallback() { }
  connectedMoveCallback() {
    console.log("Custom move-handling logic here.");
  }
}

if (!customElements.get("plc-platform-content")) {
  customElements.define("plc-platform-content", ContentComponent);
}

