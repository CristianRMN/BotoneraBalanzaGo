# App de Gestión de Turnos con UI en Go 🖥️

Proyecto desarrollado en Go con interfaz gráfica construida en Fyne, consistente en una aplicación de escritorio para la gestión de turnos y secciones.
El objetivo del proyecto es mostrar la capacidad de integrar UI dinámica, consumo de endpoints externos y una arquitectura modular en un programa de escritorio eficiente.

---

🚀 Descripción general

La aplicación abre una interfaz gráfica en la que el usuario debe completar ciertos campos obligatorios.
Una vez validados los datos, el programa realiza una llamada a un endpoint interno (disponible únicamente en la red de la empresa donde se desarrolló).

La respuesta del endpoint define la estructura de la UI:

- Si hay 2 secciones, se generan 2 botoneras.

- Si hay 5 secciones, se generan 5 botoneras.

Además, el sistema permite subir y bajar turnos en tiempo real, mostrando la información de manera clara y eficiente.

---

🛠️ Tecnologías utilizadas

- Go

- Fyne (framework para interfaces gráficas en Go)

- API REST (consumo de datos desde un servicio externo interno de la empresa)

---

🖼️ Interfaz de Usuario - UI

La UI inicial consiste en un formulario donde el usuario debe introducir los campos requeridos.

🧩 Rol dentro del sistema

📥 Validar los datos ingresados por el usuario.

🌐 Conectarse al endpoint interno para recuperar las secciones disponibles.

🔁 Construir dinámicamente la UI (botoneras) en función de la respuesta del servicio.

🧠 Funcionalidades principales
1. 📋 Validación de entradas

El formulario inicial no permite continuar si no se cubren correctamente todos los campos.

2. 🔄 Creación dinámica de botoneras

Generación de botones según el número de secciones recibidas.

3. ⬆️⬇️ Gestión de turnos

Subida y bajada de turnos en tiempo real.

Superposición eficiente de la información para mantener una experiencia fluida.

---

📚 Lo que se aprende con este proyecto

- Desarrollo de aplicaciones gráficas en Go usando Fyne.

- Consumo de endpoints REST y generación dinámica de interfaces.

- Diseño de proyectos en Go con separación de responsabilidades (UI, modelos, servicios).

- Implementación de flujos interactivos y responsivos en una aplicación de escritorio.

---

📌 Notas importantes

- El proyecto no es funcional fuera del entorno de la empresa, ya que el endpoint solo estaba disponible en su red privada.

- Este repositorio está publicado únicamente como ejemplo de estructura, arquitectura y desarrollo en Go con Fyne.


👤 Autor

Cristian Regueiro Martínez

[LinkedIn](https://www.linkedin.com/in/cristian-regueiro-mart%C3%ADnez-084187251/)

[GitHub](https://github.com/CristianRMN)
