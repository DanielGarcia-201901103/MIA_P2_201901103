import React, { useState } from "react";
import Editor from "@monaco-editor/react";
import "./CodeEditorWindow.css";
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';

import { saveAs } from 'file-saver'; // Importa la función saveAs de FileSaver.js

import { useNavigate } from 'react-router-dom';

const CodeEditorWindow = () => {
  const navigate = useNavigate();
    const [vcoman, setComandos] = useState("");
    const[outputValue, setOutputValue] = useState("");
    //const [imagen, setImagen] = useState('https://yakurefu.com/wp-content/uploads/2020/02/Chi_by_wallabby.jpg')
    const handleSubmit = async (e) => {
      setOutputValue("Cargando..." + vcoman);
      /*
        e.preventDefault();
        await fetch('http://localhost:5000/login', {
            method: 'POST',
            mode: 'cors',
            body: JSON.stringify({
                Usuario: userLogin,
                Password: passwordLogin
            }),
            headers: {
                'Access-Control-Allow-Origin': '*',
                'Content-Type': 'application/json'
            }
        })

            .then(response => response.json())
            .then(data => validar(data));
            */
        /*
        .then(function(response) {
        if (response.ok) {
            return response.json();
        }
        throw new Error('Error en la solicitud POST');
        })
        .then(function(responseData) {
            // Aquí puedes acceder a la respuesta del backend (responseData)
        })
        .catch(function(error) {
            console.log('Error:', error.message);
        });*/
    }
/*
    const validar = (data) => {
        console.log(data)
        //setImagen(data.Imagenbase64)
        if (data.data === "Administrador") {
            window.localStorage.setItem("Administrador", "201901103")
            //window.open("/admin","_self")
            navigate('/admin')
            console.log("estoy en admin")
        } else if (data.data === "SI") {
            localStorage.setItem('current', userLogin);
            window.open("/empleado", "_self");
           
        } else {
            swal({
                title: 'Credenciales Incorrectas!',
                text: 'Intenta de nuevo',
                icon: 'error',
                confirmButtonText: 'Ok'
            })
            setUsuario("")
            setPassword("")
        }
    }*/

  return (
    <>
      <div className="conjuntoBotones">
        <Form onClick={handleSubmit}>
          <InputGroup size="lg" className="custom-input-group">
            <InputGroup.Text id="inputGroup-sizing-lg">$_ </InputGroup.Text>
            <Form.Control
              placeholder="Ingresa un comando"
              type="text"
              aria-label="Large"
              aria-describedby="inputGroup-sizing-sm"
              className="custom-input"
              onChange={e => setComandos(e.target.value)}
              value={vcoman}
              autoFocus 
              required
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
            cols={50}
            placeholder="Output"
          />
        </div>
      </div>
      
      <div className="movilizando">
        <div className="spinner">
          <div></div>
          <div></div>
          <div></div>
          <div></div>
          <div></div>
          <div></div>
        </div></div>

      
      <br></br>
    </>
  );
};

export default CodeEditorWindow;