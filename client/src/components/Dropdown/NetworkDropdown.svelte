<script>
    import {onMount} from 'svelte';
    import {scale} from 'svelte/transition';

    export let user;

    let show = false; // menu state
    let menu = null; // menu wrapper DOM reference

    onMount(() => {
      const handleOutsideClick = (event) => {
        if (show && !menu.contains(event.target)) {
          show = false;
        }
      };

      // add events when element is added to the DOM
      document.addEventListener('click', handleOutsideClick, false);

      // remove events when element is removed from the DOM
      return () => {
        document.removeEventListener('click', handleOutsideClick, false);
      };
    });
</script>

<div class="w-full relative z-50" bind:this={menu}>
    <button on:click={() => (show = !show)}
            class="btn btn-outline rounded-2xl no-animation w-full origin-top-right absolute right-0 px-4 py-2">
        <span class="float-left">
            <i class="fas fa-rss"></i>
            Test Network
        </span>

        <span class="float-right ml-2">
        <i class="fas fa-chevron-down"></i>
        </span>
    </button>
    {#if show}
        <div class="w-full">
            <div in:scale={{duration: 10, start: 0.95}}
                 class="w-full origin-top-right absolute right-0 bg-primary text-primary-content rounded-2xl">
                <button on:click={() => (show = !show)}
                        class="btn btn-ghost w-full px-4 py-2 ">
                    Test Network
                </button>
                <button on:click={() => (show = !show)}
                        class="btn btn-ghost w-full px-4 py-2">
                    Main Network
                </button>
            </div>
        </div>
    {/if}
</div>