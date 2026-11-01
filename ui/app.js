// ==========================================
// LÓGICA DO FRONTEND (Go WebUI Bridge)
// ==========================================

document.addEventListener('DOMContentLoaded', () => {
    const inputNome = document.getElementById('nome');
    if (inputNome) inputNome.focus();
});

async function executarBusca() {
    const inputNome = document.getElementById('nome');
    const selectPlataforma = document.getElementById('plataforma');
    const btnSubmit = document.getElementById('btn-submit');
    const btnText = document.getElementById('btn-text');
    const btnSpinner = document.getElementById('btn-spinner');

    const statusBox = document.getElementById('status-box');
    const placeholderView = document.getElementById('placeholder-view');
    const bestDealCard = document.getElementById('best-deal-card');
    const tableCard = document.getElementById('table-card');

    const nome = inputNome.value.trim();
    const plataforma = selectPlataforma.value.trim();

    if (!nome || !plataforma) return;

    // Estado de carregamento na interface
    btnSubmit.disabled = true;
    btnText.textContent = "Pesquisando...";
    btnSpinner.style.display = "inline-block";

    statusBox.style.display = "none";
    bestDealCard.style.display = "none";
    tableCard.style.display = "none";
    placeholderView.style.display = "none";

    try {
        // Chamada à função Go exportada via WebUI
        let rawResponse;
        if (typeof buscarJogo === 'function') {
            rawResponse = await buscarJogo(nome, plataforma);
        } else if (typeof webui !== 'undefined' && typeof webui.call === 'function') {
            rawResponse = await webui.call('buscarJogo', nome, plataforma);
        } else {
            throw new Error("Comunicação WebUI com o backend Go não está disponível.");
        }

        const data = typeof rawResponse === 'string' ? JSON.parse(rawResponse) : rawResponse;

        const avisosHtml = formatAvisos(data && data.avisos);

        if (!data || !data.encontrado) {
            statusBox.className = "status-msg error";
            const mensagem = data && data.mensagem ? data.mensagem : "Não foi possível encontrar nenhum jogo com essas especificações.";
            statusBox.innerHTML = `<span>${avisosHtml ? avisosHtml + '<br>' : ''}⚠️ ${escapeHtml(mensagem)}</span>`;
            statusBox.style.display = "flex";
        } else {
            if (avisosHtml) {
                statusBox.className = "status-msg warning";
                statusBox.innerHTML = `<span>${avisosHtml}</span>`;
                statusBox.style.display = "flex";
            } else {
                statusBox.style.display = "none";
            }

            // 1. Preenche e exibe o card da MELHOR OFERTA (Menor Valor)
            document.getElementById('deal-nome').textContent = data.melhor.Nome;
            document.getElementById('deal-plataforma').textContent = data.melhor.Plataforma;
            document.getElementById('deal-preco').textContent = formatMoney(data.melhor.Preco);
            document.getElementById('deal-loja').textContent = `Loja #${data.melhor.LojaID}`;
            bestDealCard.style.display = "block";

            // 2. Preenche e exibe a tabela com TODAS as ofertas encontradas
            const tbody = document.getElementById('table-body');
            tbody.innerHTML = '';

            const total = data.resultados ? data.resultados.length : 0;
            document.getElementById('results-count').textContent = `${total} oferta(s) encontrada(s)`;

            if (data.resultados && data.resultados.length > 0) {
                data.resultados.forEach(item => {
                    const isMelhor = item.Nome === data.melhor.Nome && 
                                     item.Plataforma === data.melhor.Plataforma && 
                                     item.Preco === data.melhor.Preco && 
                                     item.LojaID === data.melhor.LojaID;

                    const tr = document.createElement('tr');
                    tr.innerHTML = `
                        <td>
                            <strong>${escapeHtml(item.Nome)}</strong>
                            ${isMelhor ? ' <span style="font-size: 11px; color: #34d399; font-weight: 600;">(Menor Preço)</span>' : ''}
                        </td>
                        <td><span class="platform-tag">${escapeHtml(item.Plataforma)}</span></td>
                        <td class="price-cell ${isMelhor ? 'lowest' : ''}">${formatMoney(item.Preco)}</td>
                        <td><span class="source-badge">⚡ Loja #${item.LojaID}</span></td>
                    `;
                    tbody.appendChild(tr);
                });
                tableCard.style.display = "block";
            }
        }
    } catch (err) {
        statusBox.className = "status-msg error";
        statusBox.innerHTML = `<span>❌ Erro de comunicação com o backend: ${escapeHtml(String(err))}</span>`;
        statusBox.style.display = "flex";
    } finally {
        btnSubmit.disabled = false;
        btnText.textContent = "Buscar Ofertas";
        btnSpinner.style.display = "none";
    }
}

function limparCampos() {
    document.getElementById('nome').value = '';
    document.getElementById('plataforma').selectedIndex = 0;
    document.getElementById('status-box').style.display = "none";
    document.getElementById('best-deal-card').style.display = "none";
    document.getElementById('table-card').style.display = "none";
    document.getElementById('placeholder-view').style.display = "flex";
    document.getElementById('nome').focus();
}

function formatMoney(val) {
    if (typeof val !== 'number') return 'R$ 0,00';
    return 'R$ ' + val.toFixed(2).replace('.', ',');
}

function formatAvisos(avisos) {
    if (!Array.isArray(avisos) || avisos.length === 0) return '';
    return `⚠️ Aviso sobre as consultas:<br>${avisos.map(aviso => `• ${escapeHtml(aviso)}`).join('<br>')}`;
}

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

