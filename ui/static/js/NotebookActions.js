function loadNotebook(documentId, appState) {
    showDocumentContent();
    if (documentId === appState.activeNotebookId) {
        return;
    }
    saveNotebook(appState); // save the old document
    appState.contentChanged = false;
    appState.activeNotebookId = documentId;

    // Send an AJAX request to retrieve document content
    fetch(`api/notebooks/${documentId}/`)
        .then(response => response.json())
        .then(data => {
            simplemde.value(data.content);
            switchToPreviewMode(simplemde);
            appState.activeNotebookId = documentId;
            appState.activeNotebookTitle = data.title;
            document.getElementById("notebook-title").innerText = appState.activeNotebookTitle;
        })
        .catch(error => console.error('Error:', error));
}

function createNewNotebook(newNotebook, appState) {
    // Send an AJAX request to create a new document
    fetch('api/notebooks/new/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-CSRFToken': `${appState.csrfToken}`, // Include CSRF token for Django
        },
        body: JSON.stringify(newNotebook),
    })
    .then(response => response.json())
    .then(data => {
        // Handle the response from the server (e.g., success or error message)
        if (data.success) {
            closeCreateNewNotebookModal();
            // Assuming the server sends a success message
            buildDocumentTitleList(appState);
            loadNotebook(data.doc_id, appState);
        } else {
            alert("Error: " + data.error);
        }
    })
    .catch(error => console.error('Error:', error));
}

function deleteNotebook(appState) {
    fetch(`api/notebooks/delete/${appState.activeNotebookId}`, {
        method: 'DELETE',
        headers: {
                'Content-Type': 'application/json',
                'X-CSRFToken': `${appState.csrfToken}`, // Include CSRF token for Django
            },
        body: JSON.stringify({}),
    })
    .then(response => {
        if (response.ok) {
            appState.activeNotebookId = null;
            buildDocumentTitleList(appState);
            hideDocumentContent();
        }
    })
    .catch(error => console.error('Error:', error));
}

function saveNotebook(appState) {
    if (!appState.contentChanged || !appState.activeNotebookId) {
        return;
    }

    fetch(`api/notebooks/save/${appState.activeNotebookId}/`, {
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
            'X-CSRFToken': `${appState.csrfToken}`,
        },
        body: JSON.stringify({ content: simplemde.value() }),
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            appState.contentChanged = false;
        }
        else {
            alert("Failed to save document");
        }
    });
}

function buildDocumentTitleList(appState, text = null) {
    let endpoint = `api/notebooks/titles/`;
    if (text !== null && text !== '') {
        endpoint = endpoint + `?text=${text}`;
    }
    fetch(`${endpoint}`)
    .then(response => response.json())
    .then(data => {
        buildNotebooksAccordion(data.notebooks, appState);
    })
    .catch(error => console.error('Error:', error));
}

function buildDocumentTitleHtmlList(notebooks) {
    const listContainer = document.getElementById('document-title-list');
    listContainer.innerHTML = '';
    notebooks.forEach(notebook => {
        const listItem = document.createElement('li');
        const strong = document.createElement('strong');
        strong.textContent = notebook.title;
        const par = document.createElement('p');
        par.textContent = "tag1;tag2;tag3;"
        listItem.appendChild(strong);

        listItem.onclick = function () {
            loadNotebook(notebook.id, '{{ csrf_token }}', appState);
        }
        listContainer.appendChild(listItem);
    });
}



function adaptNotebookToAccordionFormat(notebooks) {
    const topicToNotebooks = {};
    notebooks.forEach(notebook => {
        if (!(notebook.topic in topicToNotebooks)) {
            topicToNotebooks[notebook.topic] = [];
        }
        topicToNotebooks[notebook.topic].push({"id":notebook.id, "title": notebook.title})
    });

    const result = [];
    for (let topic in topicToNotebooks) {
        result.push({"topic":topic, "notebooks":topicToNotebooks[topic]});
    }

    return result;
}

function buildNotebooksAccordion(notebooks, appState) {
    let adaptedNotebooks = adaptNotebookToAccordionFormat(notebooks);
    const accordion = document.getElementById('accordion');
    accordion.innerHTML = ''; // Clear existing accordion items

    adaptedNotebooks.forEach(notebook => {
        // Create accordion item
        const header = document.createElement('div');
        header.className = 'accordion-header';
        header.textContent = notebook.topic;

        const content = document.createElement('div');
        content.className = 'accordion-content';

        // Populate content
        notebook.notebooks.forEach(doc => {
            const listItem = document.createElement('li');
            const strong = document.createElement('strong');
            strong.textContent = doc.title;
            listItem.appendChild(strong);

            listItem.onclick = function () {
                loadNotebook(doc.id, appState);
            }

            content.style.display = 'none';
            content.appendChild(listItem);
        });

        // Append to accordion
        accordion.appendChild(header);
        accordion.appendChild(content);

        // Event Listener for header
        header.addEventListener('click', function() {
            // Toggle this content
            content.style.display = content.style.display === 'none' ? 'block' : 'none';

            // Close other contents
            document.querySelectorAll('.accordion-content').forEach(otherContent => {
                if (otherContent !== content) {
                    otherContent.style.display = 'none';
                }
            });
        });
    });
}

function clearSearchBar() {
    document.getElementById('search-input').value = '';
}

function closeCreateNewNotebookModal() {
    addNewNotebookModal.style.display = "none";
}

function showCreateNewNotebookModal() {
    document.getElementById('add-new-notebook-title').value = '';
    document.getElementById('editable-topic-select').value = '';

    const optionsTopicSelect = document.getElementById('options-topic-select');
    while (optionsTopicSelect.firstChild) {
        optionsTopicSelect.removeChild(optionsTopicSelect.firstChild);
    }

    fetch(`api/notebooks/topics/`)
        .then(response => response.json())
        .then(data => {
            data.topics.forEach(topic => {
                const option = document.createElement('option');
                option.value = topic;
                optionsTopicSelect.appendChild(option);
            });
            addNewNotebookModal.style.display = "block";
        })
        .catch(error => {
            addNewNotebookModal.style.display = "block";
            console.error('Error:', error);
        });
}
