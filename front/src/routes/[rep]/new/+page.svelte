<script>
    import { Carta, MarkdownEditor } from 'carta-md';
    export let data;


    // Component default theme
    import 'carta-md/default.css';
    const carta = new Carta({
        // Remember to use a sanitizer to prevent XSS attacks
        // sanitizer: mySanitizer
    });
    let value = "";

    let url = '/api/upl';
    async function Read() {
        let formData = new FormData();
        formData.append('file_name', "README.md");
        formData.append('file', value);
        formData.append('repo', data.rep);

        const response = await fetch(url, {
            method: 'POST',
            body: formData,
        });
        console.log(response.status)
    }
</script>

<h1 class="my-h1">{data.rep}</h1>
<MarkdownEditor {carta} bind:value />
<button on:click={Read}>Read</button>

<style>
    /* Or in global stylesheet */
    /* Set your monospace font (Required to have the editor working correctly!) */
    :global(.carta-font-code) {
        font-family: '...', monospace;
        font-size: 1.1rem;
    }

    .my-h1 {
        color: #ff3e00;
        text-transform: uppercase;
        font-size: 4em;
        font-weight: 100;
    }
</style>