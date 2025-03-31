<script lang="ts">
  export let history = [];
  
  // Secret message function - when clicking 4 history items in a specific order
  let clickSequence = [];
  
  function handleHistoryClick(index) {
    clickSequence.push(index);
    if (clickSequence.length > 4) {
      clickSequence.shift();
    }
    
    // Check for pattern 3-1-3-7
    if (clickSequence.join('-') === '3-1-3-7' || 
        clickSequence.toString() === '3,1,3,7') {
      // Add secret message
      const secretEvent = new CustomEvent('secretfound', { 
        detail: { 
          message: 'Secret message from Maks0101aps: Thank you for exploring my calculator!',
          url: 'https://github.com/Maks0101aps'
        } 
      });
      document.dispatchEvent(secretEvent);
      
      // Reset sequence
      clickSequence = [];
    }
  }
</script>

<div class="history-panel">
  <div class="history-title">
    <span class="history-icon">📜</span>
    <span>Calculation History</span>
  </div>
  <div class="history-chat">
    {#each history as item, index}
      <div class="history-item" on:click={() => handleHistoryClick(index)}>
        <div class="history-expression">{item.expression} = </div>
        <div class="history-result">{item.result}</div>
        <div class="history-actions">
          <button class="history-action" title="Copy" on:click={() => navigator.clipboard.writeText(item.result)}>
            📋
          </button>
          <button class="history-action" title="Use" on:click={() => {
            // Event for parent component to handle
            dispatchEvent(new CustomEvent('usehistory', { detail: { value: item.result } }));
          }}>
            ↩️
          </button>
        </div>
      </div>
    {/each}
    <div class="history-footer">
      <span class="history-signature">Powered by UAH</span>
    </div>
  </div>
</div>

<style>
  .history-panel {
    margin-bottom: 20px;
    border: 3px inset #a08070;
    border-radius: 8px;
    background-color: #e0d0c0;
    box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.1);
    overflow: hidden;
  }
  
  .history-title {
    display: flex;
    align-items: center;
    background: linear-gradient(to right, #d0c0b0, #b8a898);
    padding: 5px 10px;
    font-size: 0.9rem;
    color: #5a3921;
    font-weight: bold;
    border-bottom: 2px solid #a08070;
    text-shadow: 1px 1px 0 rgba(255, 255, 255, 0.5);
  }
  
  .history-icon {
    margin-right: 8px;
    font-size: 0.8rem;
  }
  
  .history-chat {
    max-height: 150px;
    overflow-y: auto;
    padding: 8px;
    font-family: 'VT323', 'Courier New', monospace;
    background-color: #f0e8d8;
  }
  
  .history-item {
    display: flex;
    align-items: center;
    padding: 5px;
    margin-bottom: 5px;
    border-radius: 4px;
    background-color: rgba(255, 255, 255, 0.5);
    border: 1px solid #d0c0a0;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .history-item:hover {
    background-color: rgba(255, 255, 255, 0.8);
  }
  
  .history-expression {
    opacity: 0.7;
    font-size: 12px;
    margin-right: 5px;
    font-style: italic;
  }
  
  .history-result {
    font-weight: bold;
    font-size: 14px;
    margin-right: auto;
    color: #5a3921;
  }
  
  .history-actions {
    display: flex;
    gap: 3px;
  }
  
  .history-action {
    background: none;
    border: none;
    font-size: 12px;
    cursor: pointer;
    opacity: 0.6;
    transition: all 0.2s;
    border-radius: 3px;
  }
  
  .history-action:hover {
    opacity: 1;
    background-color: rgba(255, 255, 255, 0.7);
  }
  
  .history-footer {
    font-size: 0.6rem;
    text-align: right;
    padding-top: 5px;
    color: #a08070;
    opacity: 0.6;
  }
  
  .history-signature {
    font-style: italic;
  }
  
  .history-signature:hover {
    text-decoration: underline;
    cursor: help;
  }
</style> 