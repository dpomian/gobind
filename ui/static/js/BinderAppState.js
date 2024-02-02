class BinderAppState {
    constructor() {
        this.activeNotebookId = null;
        this.activeNotebookTitle = null;
        this.contentChanged = false;
        this.csrfToken = null;
        this.auth = null;
    }
}

appState = new BinderAppState();