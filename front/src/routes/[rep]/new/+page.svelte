<script>
    export let data;

    import { Carta, MarkdownEditor } from 'carta-md';
    // Component default theme
    import 'carta-md/default.css';

    // const carta = new Carta();
    // Remember to use a sanitizer to prevent XSS attacks
    // sanitizer: mySanitizer

    const carta = new Carta({
        // Remember to use a sanitizer to prevent XSS attacks
        // sanitizer: mySanitizer
    });

    function Read() {
        console.log("huh")
        console.log(carta.renderer.container.contentEditable.valueOf().)
        carta.cartaListeners.forEach(s => {
            console.log(s)
        })
    }
</script>

<h1 class="my-h1">{data.rep}</h1>
<MarkdownEditor {carta} />
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