class NavigationComponent extends HTMLElement {
  static observedAttributes = [];

  constructor() {
    super();
  }

  attributeChangedCallback(name, oldValue, newValue) {
    console.log(`Attribute ${name} has changed.`);
  }
  connectedCallback() {
    console.log("Navigation connected");
  }
  disconnectedCallback() { }
  connectedMoveCallback() {
    console.log("Custom move-handling logic here.");
  }
}

if (!customElements.get("plc-platform-navigation")) {
  customElements.define("plc-platform-navigation", NavigationComponent);
}

