let formData = new FormData();

function showScreenById(screenId) {
    const screenToShow = document.getElementById(screenId);

    if (screenToShow) {
        document.querySelectorAll('[data-type="screen"]').forEach(screen => {
            screen.classList.add('display-none');
        });
        screenToShow.classList.remove('display-none');
    }
}

function updateLoadedFilesCount() {
    const loadedFilesCount = formData.getAll('files').length + document.getElementById('files').files.length;
    document.getElementById('loaded-files-count').innerText = `${loadedFilesCount} files loaded.`;
}

function uploadFiles() {
    const filesInput = document.getElementById('files');
    [...filesInput.files].map(f => formData.append('files', f));

    if (formData.getAll('files').length === 0) {
        return;
    }

    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/upload", true);
    xhr.onload = () => {
        if (xhr.status === 200) {
            formData = new FormData();
            document.getElementById("upload-files-form").reset();
            showScreenById('success-screen');
        } else {
            console.log("File upload failed");
        }
    };
    showScreenById('uploading-screen');
    xhr.send(formData);
}

async function showExistingUploadsPage() {
    showScreenById('existing-uploads-screen')
    await loadExistingUploads();
}

async function loadExistingUploads() {
    const existingUploadsListElement = document.getElementById('existing-uploads-list');
    try {
        const respose = await fetch('/existingUploads');
        const existingUploads = await respose.json();

        const itemNodes = existingUploads.map((upload) => {
            const { name, path } = upload;
            const a = document.createElement('a', );
            a.href = path;
            a.text = name;
            return a;
        })

        if (itemNodes.length === 0) {
            throw new Error("no existing uploads");
        }

        existingUploadsListElement.replaceChildren(...itemNodes);
    }
    catch (e) {
        console.error(e);
        existingUploadsListElement.replaceChildren("No existing uploads found.")
    }
}

function handleDrop(e) {
    e.preventDefault();

    const dropItems = e.dataTransfer.items ?? e.dataTransfer.files;
    const files = [...dropItems].map((item) => {
        if (item.kind === 'file') {
            return item.getAsFile();
        }
        return null;
    }).filter(Boolean);
    files.map(f => formData.append('files', f));
    updateLoadedFilesCount();
}

function handleDragOver(e) {
    e.preventDefault();
}
