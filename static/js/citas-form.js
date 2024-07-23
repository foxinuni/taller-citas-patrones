const form = document.getElementById('citas-form');
if (!form) {
    throw new Error('Form was not found');
}

form.addEventListener('submit', async (event) => {
    event.preventDefault();

    const formData = new FormData(form);
    const data = Object.fromEntries(formData.entries());

    const response = await fetch(form.action, {
        method: form.method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            persona: {
                nombre: data.nombre,
                apellido: data.apellido,
                cedula: data.cedula,
                edad: parseInt(data.edad),
            },
            fecha: data.fecha + ':05Z',
        }),
    });
    
    if (response.ok) {
        window.location.href = '/citas';
    } else {
        const error = await response.json();
        alert(error);
    }
});