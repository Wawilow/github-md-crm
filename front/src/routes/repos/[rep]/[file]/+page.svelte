<svelte:options immutable={false} />

<script lang="ts">
    import { onMount } from "svelte";
    export let data;
    import { Carta, MarkdownEditor } from 'carta-md';
    import 'carta-md/default.css';

    const carta = new Carta({
        // Remember to use a sanitizer to prevent XSS attacks
        // sanitizer: mySanitizer
    });
    let value = "";

    const url: string = '/api/file?';
    async function getReps(): Promise<Array<any>> {
        try {
            const response = await fetch(url  + new URLSearchParams({repo: data.rep, path: data.file}).toString());
            if (response.status !== 200) {
                console.log("new file");
            } else {
            value = await response.text();
        }
        } catch (error) {
            console.error('Error:', error);
        }
        return []
    }

    onMount(async () => {
        await getReps()
    })

    let url2 = '/api/upl';
    async function SendFile() {
        let formData = new FormData();
        formData.append('file', data.file);
        formData.append('content', value);
        formData.append('repo', data.rep);
        console.log(formData)
        const response = await fetch(url2, {
            method: 'POST',
            body: formData,
        });
        console.log(response.status)
        window.location.reload()
    }
</script>

<main>
    <h1 class="my-h1">File {data.rep}/{data.file}</h1>
    <MarkdownEditor {carta} bind:value />
    <button on:click={SendFile}>Update file</button>
</main>

<style>

    h1 {
        color: #ff3e00;
        text-transform: uppercase;
        font-size: 4em;
        font-weight: 100;
    }

    table {
        color: #ff3e00;
        margin-left: auto;
        margin-right: auto;
        font-size: 1.5em;
        font-weight: 100;
    }

    @media (min-width: 640px) {
        main {
            max-width: none;
        }
    }
</style>