import React, { useState } from "react";
import Editor from "@monaco-editor/react";
import "./CodeEditorWindow.css";
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';

import { saveAs } from 'file-saver'; // Importa la función saveAs de FileSaver.js


const CodeEditorWindow = () => {
  const [inputValue, setInputValue] = useState("");
  const [outputValue, setOutputValue] = useState("");

  const handleEditorChange = (value) => {
    setInputValue(value);
  };

  const handleRunCode = async (e) => {
    setOutputValue("");
    if (inputValue != "") {
      e.preventDefault();
      await fetch('http://localhost:4000/analizar', {
        method: 'POST',
        mode: 'cors',
        body: JSON.stringify({ texto: inputValue }),
        headers: {
          'Access-Control-Allow-Origin': '*',
          'Content-Type': 'application/json'
        }
      }).then(response => response.json()).then(data => validar(data));

      // Aquí puedes implementar la lógica para ejecutar el código y actualizar la salida
      // Ejemplo: establecer la salida como el valor de entrada
    }

  };
  const validar = (data) => {
    let texts = "";
    for (let i = 0; i < data.salida.length; i++) {
      texts += data.salida[i];
    }
    setOutputValue(texts);
  };

  const handleOpenFile = async (e) => {
    const file = e.target.files[0]; // Obtiene el archivo seleccionado por el usuario
    const fileContent = await readFileContent(file); // Lee el contenido del archivo
    setInputValue(fileContent); // Actualiza el estado con el contenido del archivo
  };

  const readFileContent = (file) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => resolve(e.target.result);
      reader.onerror = (e) => reject(e);
      reader.readAsText(file);
    });
  };

  const handleSaveFile = () => {
    const blob = new Blob([inputValue], { type: 'text/plain;charset=utf-8' }); // Crea un Blob con el contenido del editor
    saveAs(blob, "archivo.sc"); // Guarda el archivo con el nombre "archivo.txt"
  };

  const handleCreateFile = () => {
    const emptyContent = ""; // Contenido en blanco para el archivo
    setInputValue(emptyContent)

    const blob = new Blob([emptyContent], { type: 'text/plain;charset=utf-8' }); // Crea un Blob con el contenido en blanco
    saveAs(blob, "nuevoArchivo.sc"); // Guarda el archivo con la extensión .sc y el nombre "nuevoArchivo"
  };
  const handleOpenReportErrors = async (e) => {
    e.preventDefault();
    await fetch('http://localhost:4000/erroresT', {
      method: 'GET',
      mode: 'cors',
      headers: {
        'Access-Control-Allow-Origin': '*',
        'Content-Type': 'application/json'
      }
    }).then(response => response.json()).then(data => validare(data));

    // Aquí puedes implementar la lógica para ejecutar el código y actualizar la salida
    // Ejemplo: establecer la salida como el valor de entrada
  };
  const validare = (data) => {
    console.log(data.salida);
  };
  return (
    <>
      <div className="conjuntoBotones">
        <Form>
          <InputGroup size="lg" className="custom-input-group">
            <InputGroup.Text id="inputGroup-sizing-lg">$_ </InputGroup.Text>
            <Form.Control
              placeholder="Ingresa un comando"
              aria-label="Large"
          aria-describedby="inputGroup-sizing-sm"
          className="custom-input" 
            />
              <Button variant="primary" type="button" className="custom-submit-button">
                Enviar
              </Button>
          </InputGroup>

        </Form>
      </div>
      <div className="code-editor-container">

        <div className="code-editor-right">
          <h2>Consola</h2>
          <textarea
            className="output-console"
            value={outputValue}
            readOnly
            rows={inputValue.split("\n").length}
            cols={50}
            placeholder="Output"
          />
        </div>
      </div>
    </>
  );
};

export default CodeEditorWindow;