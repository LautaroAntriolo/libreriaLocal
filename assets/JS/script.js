// static/js/script.js
document.addEventListener('DOMContentLoaded', function() {
    // Función para actualizar campos dependientes
    const leidoSelect = document.getElementById('leido');
    const fechaLecturaInput = document.getElementById('fecha_lectura');
    const puntajeInput = document.getElementById('puntaje');

    if (leidoSelect) {
        leidoSelect.addEventListener('change', function() {
            if (this.value === '1') { // Si está leído
                porcentajeInput.value = 100;
                
                // Si no hay fecha de lectura, poner la fecha actual
                if (fechaLecturaInput.value === '') {
                    const hoy = new Date();
                    const formattedDate = hoy.toISOString().split('T')[0];
                    fechaLecturaInput.value = formattedDate;
                }
            } else { // Si no está leído
                if (porcentajeInput.value === '100') {
                    porcentajeInput.value = 0;
                }
            }
        });
    }

    // Validación del formulario
    const form = document.querySelector('form');
    if (form) {
        form.addEventListener('submit', function(event) {
            const nombre = document.getElementById('nombre');
            const autor = document.getElementById('autor');
            const editorial = document.getElementById('editorial');
            
            if (nombre.value.trim() === '') {
                alert('El nombre del libro es obligatorio');
                event.preventDefault();
                nombre.focus();
                return false;
            }
            
            if (autor.value.trim() === '') {
                alert('El autor del libro es obligatorio');
                event.preventDefault();
                autor.focus();
                return false;
            }
            
            if (editorial.value.trim() === '') {
                alert('La editorial es obligatoria');
                event.preventDefault();
                editorial.focus();
                return false;
            }
            
            return true;
        });
    }
    
    // Actualizar porcentaje automáticamente cuando se cambia a leído
    if (porcentajeInput && leidoSelect) {
        porcentajeInput.addEventListener('change', function() {
            if (parseInt(this.value) === 100 && leidoSelect.value === '0') {
                leidoSelect.value = '1';
            } else if (parseInt(this.value) < 100 && leidoSelect.value === '1') {
                leidoSelect.value = '0';
            }
        });
    }
});