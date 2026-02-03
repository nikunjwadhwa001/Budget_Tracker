const form = document.getElementById('budgetForm');
const list = document.getElementById('list');

// Load items when page opens
async function loadTransactions() {
    const res = await fetch('/transactions');
    const data = await res.json();

    // Buffalo API returns clean JSON arrays by default
    const items = data || [];

    list.innerHTML = items.map(item => `
        <li class="${item.type}">
            <strong>${item.description}</strong>: $${item.amount}
        </li>
    `).join('');
}

// Handle form submit
form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const desc = document.getElementById('desc').value;
    const amount = parseFloat(document.getElementById('amount').value);
    const type = document.getElementById('type').value;

    await fetch('/transactions', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        // Buffalo expects snake_case keys by default
        body: JSON.stringify({ description: desc, amount: amount, type: type })
    });

    form.reset();
    loadTransactions();
});

loadTransactions();