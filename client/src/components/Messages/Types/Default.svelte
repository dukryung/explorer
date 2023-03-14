<script>
    import {_} from 'svelte-i18n';

    export let message;

    function syntaxHighlight(json) {
      json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
      return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function(match) {
        let cls = 'txMessage-default-blue-400';
        if (/^"/.test(match)) {
          if (/:$/.test(match)) {
            cls = 'txMessage-default-red-400';
          } else {
            cls = 'txMessage-default-green-400';
          }
        } else if (/true|false/.test(match)) {
          cls = 'txMessage-default-blue-400';
        } else if (/null/.test(match)) {
          cls = 'txMessage-default-pink-400';
        }
        return `<span class="${cls}">${match}</span>`;
      });
    }
</script>

<style>

</style>

<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{message['@type']}</h5></div>
                </div>
            </div>
            <div style="margin-top: 2rem; margin-bottom: 2rem; background-color: black; color:white">
                <pre style="text-align: left">
                {@html syntaxHighlight(JSON.stringify(message, undefined, '\t'))}
                </pre>
            </div>
        </div>
    </div>
</section>