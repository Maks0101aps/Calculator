<script lang="ts">
  import { onMount } from 'svelte';
  import { invoke } from '@tauri-apps/api/tauri';
  import CalculatorHistory from './CalculatorHistory.svelte';

  let displayValue = '0';
  let currentCalculation = '';
  let history = [];
  let isScientificMode = false;
  let speakResult = false;
  let isListening = false;
  let isCurrencyMode = false;
  let childMode = false;
  let selectedCurrency = 'USD';
  let exchangeRates = {
    USD: 1,
    EUR: 0.94,
    UAH: 36.74
  };
  let isLoadingRates = false;
  // Easter eggs
  let clickCount = 0;
  let secretCode = '';
  let isDeveloperMode = false;
  let rainbowMode = false;
  
  function clearDisplay() {
    displayValue = '0';
    currentCalculation = '';
  }
  
  function appendToDisplay(value) {
    if (displayValue === '0' && value !== '.') {
      displayValue = value;
    } else {
      displayValue += value;
    }
    currentCalculation += value;
  }
  
  function handleOperator(operator) {
    if (currentCalculation !== '') {
      if (['+', '-', '*', '/'].includes(currentCalculation.slice(-1))) {
        currentCalculation = currentCalculation.slice(0, -1) + operator;
      } else {
        currentCalculation += operator;
      }
      displayValue = '0';
    }
  }
  
  async function calculateResult() {
    if (currentCalculation !== '') {
      try {
        // Basic sanitization and validation
        if (!/^[0-9+\-*/().\s^√πe]+$/.test(currentCalculation)) {
          throw new Error('Invalid input');
        }
        
        // Prepare expression for evaluation
        let expression = currentCalculation
          .replace(/√/g, 'Math.sqrt')
          .replace(/π/g, 'Math.PI')
          .replace(/e/g, 'Math.E')
          .replace(/\^/g, '**');
        
        // Calculate locally for simple expressions or use backend for complex ones
        let result;
        if (isScientificMode || expression.includes('Math.') || expression.includes('**')) {
          // Call backend for scientific calculations
          result = await invoke('calculate', { expression });
        } else {
          // Simple calculations can be done locally
          // eslint-disable-next-line no-new-func
          result = Function('"use strict"; return (' + expression + ')')();
        }
        
        // Format result
        result = parseFloat(Number(result).toFixed(10)).toString();
        
        // Add to history
        history = [...history, { expression: currentCalculation, result }];
        
        // Speak result if enabled
        if (speakResult) {
          speakText(`The result is ${result}`);
        }
        
        displayValue = result;
        currentCalculation = result;
        
      } catch (error) {
        displayValue = 'Error';
        currentCalculation = '';
      }
    }
  }
  
  function toggleScientificMode() {
    isScientificMode = !isScientificMode;
    if (isCurrencyMode && isScientificMode) {
      isCurrencyMode = false;
    }
  }
  
  function toggleCurrencyMode() {
    isCurrencyMode = !isCurrencyMode;
    if (isCurrencyMode && isScientificMode) {
      isScientificMode = false;
    }
    if (isCurrencyMode) {
      // Fetch latest exchange rates when currency mode is activated
      fetchExchangeRates();
    }
  }
  
  function toggleChildMode() {
    childMode = !childMode;
    if (childMode) {
      // Initialize child mode with friendly display
      displayValue = 'Hi there!';
      setTimeout(() => {
        displayValue = '0';
      }, 2000);
    }
  }
  
  async function fetchExchangeRates() {
    try {
      isLoadingRates = true;
      displayValue = 'Loading rates...';
      
      // Используем Exchangerate API - это открытый API без ключа
      const response = await fetch('https://api.exchangerate.host/latest?base=USD');
      const data = await response.json();
      
      if (data && data.rates) {
        // Обновляем курсы валют из полученных данных
        exchangeRates = {
          USD: 1, // базовая валюта всегда 1
          EUR: data.rates.EUR || exchangeRates.EUR, // сохраняем текущий курс, если новый не получен
          UAH: data.rates.UAH || exchangeRates.UAH
        };
        
        // Отображаем текущее значение с префиксом валюты
        const numericValue = displayValue === 'Loading rates...' ? '0' : displayValue;
        displayValue = `${selectedCurrency}: ${numericValue}`;
      } else {
        throw new Error('Failed to fetch exchange rates');
      }
    } catch (error) {
      console.error('Failed to fetch exchange rates', error);
      // В случае ошибки просто используем статические курсы
      displayValue = `${selectedCurrency}: ${displayValue}`;
    } finally {
      isLoadingRates = false;
    }
  }
  
  function convertCurrency(targetCurrency) {
    if (displayValue === '0' || displayValue === 'Error' || displayValue === 'Loading rates...') return;
    
    try {
      // Extract numeric value if display has currency prefix
      let numericValue = displayValue;
      let sourceCurrency = selectedCurrency;
      
      if (displayValue.includes(':')) {
        const parts = displayValue.split(':');
        sourceCurrency = parts[0].trim();
        numericValue = parts[1].trim();
      }
      
      // Ensure we have a valid number
      const numValue = parseFloat(numericValue);
      if (isNaN(numValue)) {
        displayValue = `${targetCurrency}: 0`;
        selectedCurrency = targetCurrency;
        return;
      }
      
      // Convert from current currency to target currency
      const valueInUSD = numValue / exchangeRates[sourceCurrency];
      const convertedValue = valueInUSD * exchangeRates[targetCurrency];
      
      // Update display with new currency and value
      selectedCurrency = targetCurrency;
      displayValue = `${targetCurrency}: ${convertedValue.toFixed(2)}`;
      currentCalculation = convertedValue.toString();
      
      // Add to history
      history = [...history, { 
        expression: `Convert ${numValue} ${sourceCurrency} to ${targetCurrency}`, 
        result: `${targetCurrency}: ${convertedValue.toFixed(2)}`
      }];
    } catch (error) {
      console.error('Currency conversion error:', error);
      displayValue = `${targetCurrency}: 0`;
      selectedCurrency = targetCurrency;
    }
  }
  
  function generateRandomNumber() {
    const randomNum = Math.floor(Math.random() * 100) + 1;
    displayValue = randomNum.toString();
    currentCalculation = randomNum.toString();
  }
  
  function handleSpecialFunction(func) {
    switch(func) {
      case 'sin':
        currentCalculation = `Math.sin(${displayValue})`;
        calculateResult();
        break;
      case 'cos':
        currentCalculation = `Math.cos(${displayValue})`;
        calculateResult();
        break;
      case 'tan':
        currentCalculation = `Math.tan(${displayValue})`;
        calculateResult();
        break;
      case 'log':
        currentCalculation = `Math.log(${displayValue})`;
        calculateResult();
        break;
      case 'random':
        generateRandomNumber();
        break;
      default:
        break;
    }
  }
  
  function toggleVoiceInput() {
    isListening = !isListening;
    if (isListening) {
      startVoiceRecognition();
    } else {
      stopVoiceRecognition();
    }
  }
  
  function toggleSpeakResult() {
    speakResult = !speakResult;
  }
  
  function speakText(text) {
    if ('speechSynthesis' in window) {
      const utterance = new SpeechSynthesisUtterance(text);
      window.speechSynthesis.speak(utterance);
    }
  }
  
  function startVoiceRecognition() {
    if ('webkitSpeechRecognition' in window) {
      const recognition = new webkitSpeechRecognition();
      recognition.continuous = false;
      recognition.interimResults = false;
      recognition.lang = 'en-US';
      
      recognition.onresult = (event) => {
        const transcript = event.results[0][0].transcript;
        processVoiceInput(transcript);
      };
      
      recognition.onerror = (event) => {
        console.error('Speech recognition error', event.error);
        isListening = false;
      };
      
      recognition.onend = () => {
        isListening = false;
      };
      
      recognition.start();
    } else {
      console.error('Speech recognition not supported');
      isListening = false;
    }
  }
  
  function stopVoiceRecognition() {
    if ('webkitSpeechRecognition' in window) {
      // Stop the recognition
    }
  }
  
  function processVoiceInput(transcript) {
    // Clean up transcript and process commands
    const cleaned = transcript.toLowerCase().trim();
    
    if (cleaned.includes('clear')) {
      clearDisplay();
    } else if (cleaned.includes('calculate') || cleaned.includes('equals')) {
      calculateResult();
    } else if (cleaned.includes('random number')) {
      generateRandomNumber();
    } else if (cleaned.includes('scientific mode')) {
      toggleScientificMode();
    } else if (cleaned.includes('currency mode')) {
      toggleCurrencyMode();
    } else if (cleaned.match(/\d+(\s*(plus|minus|times|divided by)\s*\d+)+/)) {
      // Process arithmetic expressions
      const mappedExpression = cleaned
        .replace(/plus/g, '+')
        .replace(/minus/g, '-')
        .replace(/times/g, '*')
        .replace(/divided by/g, '/');
      
      // Extract numbers and operations
      const matches = mappedExpression.match(/\d+|[+\-*/]/g);
      if (matches) {
        currentCalculation = matches.join('');
        calculateResult();
      }
    }
  }
  
  // Handle keyboard input
  function handleKeydown(event) {
    if (event.key >= '0' && event.key <= '9') {
      appendToDisplay(event.key);
    } else if (event.key === '.') {
      appendToDisplay('.');
    } else if (['+', '-', '*', '/'].includes(event.key)) {
      handleOperator(event.key);
    } else if (event.key === 'Enter' || event.key === '=') {
      calculateResult();
    } else if (event.key === 'Escape' || event.key === 'c' || event.key === 'C') {
      clearDisplay();
    }
    
    // Easter egg: Developer mode with Ctrl+Shift+D
    if (event.ctrlKey && event.shiftKey && event.key === 'D') {
      toggleDeveloperMode();
    }
    
    // Track sequence for rainbow easter egg
    trackSecretCode(event.key);
  }
  
  function toggleDeveloperMode() {
    isDeveloperMode = !isDeveloperMode;
    if (isDeveloperMode) {
      history = [...history, { 
        expression: 'Developer Mode Activated', 
        result: 'Created by Maks0101aps' 
      }];
    }
  }
  
  function trackSecretCode(key) {
    secretCode += key;
    if (secretCode.length > 10) {
      secretCode = secretCode.substring(1);
    }
    
    // Check for "rainbow" sequence
    if (secretCode.includes('rainbow')) {
      rainbowMode = true;
      history = [...history, { 
        expression: '🌈 Rainbow Mode Activated', 
        result: 'Colorful greetings from Maks0101aps!' 
      }];
      
      // Reset the code after activation
      secretCode = '';
    }
    
    // Check for "1337" sequence
    if (secretCode.includes('1337')) {
      history = [...history, { 
        expression: 'Hacker Mode Activated', 
        result: 'https://github.com/Maks0101aps' 
      }];
    }
  }
  
  function handleCalcClick() {
    clickCount++;
    if (clickCount === 5) {
      // Show serial number after 5 clicks
      history = [...history, { 
        expression: 'Serial Number', 
        result: 'Calculator-Maks0101aps-UA-2025' 
      }];
      clickCount = 0;
    }
  }
  
  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    return () => {
      window.removeEventListener('keydown', handleKeydown);
    };
  });
</script>

<div class="vintage-calculator" class:rainbow-theme={rainbowMode} on:click={handleCalcClick}>
  <div class="calculator-branding">
    <div class="calculator-logo">CALC-O-MATIC</div>
    <div class="calculator-model">Model 3000</div>
    {#if isDeveloperMode}
      <div class="dev-signature">by Maks0101aps</div>
    {/if}
  </div>
  
  <div class="solar-panel">
    <div class="solar-cell"></div>
    <div class="solar-cell"></div>
    <div class="solar-cell"></div>
    <div class="solar-cell"></div>
    <div class="solar-cell"></div>
  </div>
  
  <div class="calculator-display">
    <div class="display-glare"></div>
    <div class="expression">{currentCalculation || '0'}</div>
    <div class="result">
      {#if isLoadingRates}
        <div class="loading-indicator">Loading...</div>
      {:else}
        {displayValue}
      {/if}
    </div>
    <div class="display-indicator"></div>
  </div>
  
  {#if history.length > 0}
    <CalculatorHistory {history} />
  {/if}
  
  <div class="mode-toggles">
    <button class="mode-button" class:active={isScientificMode} on:click={toggleScientificMode}>
      Scientific
    </button>
    <button class="mode-button" class:active={isCurrencyMode} on:click={toggleCurrencyMode}>
      Currency
    </button>
    <button class="mode-button" class:active={childMode} on:click={toggleChildMode}>
      Child Mode
    </button>
    <button class="mode-button" class:active={speakResult} on:click={toggleSpeakResult}>
      Speak
    </button>
    <button class="mode-button" class:active={isListening} on:click={toggleVoiceInput}>
      Voice
    </button>
  </div>
  
  <div class="calculator-buttons">
    {#if isScientificMode}
      <!-- Scientific calculator buttons -->
      <button class="calculator-button function-btn" on:click={() => handleSpecialFunction('sin')}>sin</button>
      <button class="calculator-button function-btn" on:click={() => handleSpecialFunction('cos')}>cos</button>
      <button class="calculator-button function-btn" on:click={() => handleSpecialFunction('tan')}>tan</button>
      <button class="calculator-button function-btn" on:click={() => handleSpecialFunction('log')}>log</button>
      <button class="calculator-button function-btn" on:click={() => appendToDisplay('π')}>π</button>
      <button class="calculator-button function-btn" on:click={() => appendToDisplay('e')}>e</button>
      <button class="calculator-button function-btn" on:click={() => appendToDisplay('^')}>^</button>
      <button class="calculator-button function-btn" on:click={() => appendToDisplay('√')}>√</button>
    {/if}
    
    {#if isCurrencyMode}
      <!-- Currency converter buttons -->
      <div class="currency-header">Конвертация валют</div>
      <button class="calculator-button currency-btn" on:click={() => convertCurrency('USD')}>USD</button>
      <button class="calculator-button currency-btn" on:click={() => convertCurrency('EUR')}>EUR</button>
      <button class="calculator-button currency-btn" on:click={() => convertCurrency('UAH')}>UAH</button>
      <div class="currency-info">Курс: {selectedCurrency} = ${exchangeRates[selectedCurrency]}</div>
    {/if}
    
    <!-- Standard calculator buttons -->
    <button class="calculator-button clear-btn" on:click={clearDisplay}>C</button>
    <button class="calculator-button" on:click={() => appendToDisplay('(')}>(</button>
    <button class="calculator-button" on:click={() => appendToDisplay(')')}>)</button>
    <button class="calculator-button operator" on:click={() => handleOperator('/')}>/</button>
    
    <button class="calculator-button number-btn" on:click={() => appendToDisplay('7')}>7</button>
    <button class="calculator-button number-btn" on:click={() => appendToDisplay('8')}>8</button>
    <button class="calculator-button number-btn" on:click={() => appendToDisplay('9')}>9</button>
    <button class="calculator-button operator" on:click={() => handleOperator('*')}>*</button>
    
    <button class="calculator-button number-btn" on:click={() => appendToDisplay('4')}>4</button>
    <button class="calculator-button number-btn" on:click={() => appendToDisplay('5')}>5</button>
    <button class="calculator-button number-btn" on:click={() => appendToDisplay('6')}>6</button>
    <button class="calculator-button operator" on:click={() => handleOperator('-')}>-</button>
    
    <button class="calculator-button number-btn" on:click={() => appendToDisplay('1')}>1</button>
    <button class="calculator-button number-btn" on:click={() => appendToDisplay('2')}>2</button>
    <button class="calculator-button number-btn" on:click={() => appendToDisplay('3')}>3</button>
    <button class="calculator-button operator" on:click={() => handleOperator('+')}>+</button>
    
    <button class="calculator-button" on:click={() => handleSpecialFunction('random')}>?</button>
    <button class="calculator-button number-btn" on:click={() => appendToDisplay('0')}>0</button>
    <button class="calculator-button" on:click={() => appendToDisplay('.')}>.</button>
    <button class="calculator-button equals" on:click={calculateResult}>=</button>
  </div>
  
  <div class="calculator-footer">
    <div class="vintage-screws">
      <div class="screw"></div>
      <div class="screw"></div>
      <div class="screw"></div>
      <div class="screw"></div>
    </div>
  </div>
</div>

<style>
  /* These styles add to the global styles defined in app.css */
  .calculator-branding {
    text-align: center;
    margin-bottom: 15px;
    color: #5a3921;
  }
  
  .calculator-logo {
    font-size: 1.2rem;
    font-weight: bold;
    text-shadow: 1px 1px 0 rgba(255, 255, 255, 0.7);
    letter-spacing: 1px;
  }
  
  .calculator-model {
    font-size: 0.7rem;
    opacity: 0.8;
  }
  
  .display-glare {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, rgba(255, 255, 255, 0.3) 0%, rgba(255, 255, 255, 0) 50%);
    pointer-events: none;
    z-index: 1;
  }
  
  .display-indicator {
    position: absolute;
    right: 10px;
    bottom: 10px;
    width: 8px;
    height: 8px;
    background-color: #68a768;
    border-radius: 50%;
    box-shadow: 0 0 5px #68a768;
    animation: blink 2s infinite;
  }
  
  .expression {
    font-size: 14px;
    opacity: 0.7;
    overflow-x: auto;
    height: 20px;
    margin-bottom: 5px;
    z-index: 2;
    position: relative;
  }
  
  .result {
    font-size: 32px;
    font-weight: bold;
    overflow-x: auto;
    z-index: 2;
    position: relative;
  }
  
  .solar-panel {
    display: flex;
    justify-content: center;
    gap: 4px;
    margin-bottom: 15px;
    padding: 5px;
    background-color: #333;
    border-radius: 3px;
    box-shadow: inset 0 0 3px rgba(0, 0, 0, 0.5);
  }
  
  .solar-cell {
    width: 30px;
    height: 15px;
    background: linear-gradient(to bottom, #264040 0%, #366060 100%);
    border: 1px solid #222;
  }
  
  .function-btn {
    background-color: #d8c0a0;
    font-size: 1.2rem;
  }
  
  .currency-btn {
    background-color: #c0d8c0;
    font-size: 1rem;
    font-weight: bold;
    color: #2c4c2c;
  }
  
  .currency-header {
    grid-column: span 4;
    text-align: center;
    font-size: 1rem;
    font-weight: bold;
    color: #5a3921;
    margin-bottom: 5px;
    background-color: #e0c0a0;
    padding: 5px;
    border-radius: 5px;
    border: 2px solid #a08060;
  }
  
  .currency-info {
    grid-column: span 4;
    text-align: center;
    font-size: 0.8rem;
    color: #5a3921;
    margin-top: 5px;
    background-color: #f0e8d8;
    padding: 5px;
    border-radius: 5px;
    border: 1px dashed #a08060;
  }
  
  .clear-btn {
    background-color: #e57373;
    color: #5a2929;
  }
  
  .number-btn {
    font-family: 'VT323', monospace;
    font-size: 1.7rem;
  }
  
  .calculator-footer {
    margin-top: 20px;
    text-align: center;
  }
  
  .vintage-screws {
    display: flex;
    justify-content: space-between;
    margin: 0 10px;
  }
  
  .screw {
    width: 12px;
    height: 12px;
    background: radial-gradient(circle at 30% 30%, #aaa, #777);
    border-radius: 50%;
    position: relative;
  }
  
  .screw::after {
    content: "+";
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-size: 8px;
    color: #555;
  }
  
  .loading-indicator {
    font-size: 24px;
    color: #336633;
    opacity: 0.8;
    text-align: center;
    animation: pulse 1.5s infinite;
  }
  
  .dev-signature {
    font-size: 0.6rem;
    margin-top: 2px;
    font-style: italic;
    color: #8b5a2b;
    text-decoration: underline;
  }
  
  /* Rainbow theme easter egg */
  .rainbow-theme {
    animation: rainbow-bg 10s linear infinite;
  }
  
  @keyframes rainbow-bg {
    0% { background-color: #e8d4b8; }
    20% { background-color: #d4e8b8; }
    40% { background-color: #b8d4e8; }
    60% { background-color: #d4b8e8; }
    80% { background-color: #e8b8d4; }
    100% { background-color: #e8d4b8; }
  }
  
  @keyframes pulse {
    0% { opacity: 0.4; }
    50% { opacity: 1; }
    100% { opacity: 0.4; }
  }
  
  @keyframes blink {
    0% { opacity: 0.2; }
    50% { opacity: 1; }
    100% { opacity: 0.2; }
  }
</style> 