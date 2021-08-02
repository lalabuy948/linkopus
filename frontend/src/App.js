import './styles/App.css';
import "bootstrap/dist/css/bootstrap.min.css";
import "shards-ui/dist/css/shards.min.css"

import React, { useState, useEffect } from 'react';
import {
    Container,
    FormInput,
    Button,
    InputGroupAddon,
    InputGroupText,
    InputGroup,
    Row,
    Col,
    Collapse,
    ListGroup,
    ListGroupItem,
    ListGroupItemHeading,
    Popover,
    PopoverHeader,
    PopoverBody,
} from "shards-react";
import { CopyToClipboard } from 'react-copy-to-clipboard';

import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faLink } from '@fortawesome/free-solid-svg-icons';

import {shortLink, getLinkHash} from "./service/LinkMapService";
import {getTopLinkStats} from "./service/LinkViewService";
import {isValidURL} from "./service/UtilityService";

function App() {
  const [topLinks, setTopLinks] = useState([]);

  const [userLink, setUserLink] = useState("");
  const [urlToCopy, setUrlToCopy] = useState("");
  const [urlToCopyCopied, setUrlToCopyCopied] = useState(false);

  const [userLinkValid, setUserLinkValid] = useState(false);
  const [userLinkInvalid, setUserLinkInvalid] = useState(false);
  const [showLinkHash, setShowLinkHash] = useState(false);

  useEffect(() => {
    getTopLinkStats().then((response) => {
      if (response.data !== null) {
        setTopLinks(response.data)
      }
    }, (error) => {
        console.log(error)
    })
  }, [])

  function handleUserLinkChange(e) {
      let link = e.target.value;

      if ( !isValidURL(link) ) {
          setUserLinkInvalid(true)
          setUserLinkValid(false)
      } else {
          setUserLinkInvalid(false)
          setUserLinkValid(true)
      }

      setUserLink(link)
  }

  function handleUserLinkSubmit(link) {
      if ( isValidURL(link) ) {
          shortLink(link)
              .then((response) => {
                  setShowLinkHash(true)
                  setUrlToCopy( process.env.REACT_APP_DOMAIN_BASE_URL + "/" + response.data.link_hash)
              }, (error) => {
                  if (error.response.data.error_msg === "service: mapping already exists") {
                      getLinkHash(link)
                          .then((response) => {
                            setShowLinkHash(true)
                            setUrlToCopy( process.env.REACT_APP_DOMAIN_BASE_URL + "/" + response.data.link_hash)
                      }, (error) => {
                          setUserLinkInvalid(true)
                          setUserLinkValid(false)
                          setUrlToCopy('🤷🏻‍♀️ ups 🤷🏻‍♂️' + error);
                          setShowLinkHash(true)
                      })
                  } else {
                      setUserLinkInvalid(true)
                      setUserLinkValid(false)
                      setUrlToCopy(
                          typeof error.response !== 'undefined'
                              ? error.response.data.error_msg
                              : '🤷🏻‍♀️ ups 🤷🏻‍♂️'
                      );
                      setShowLinkHash(true)
                  }
              });
      }
  }

  function handleUserLinkSubmitOnKey(event) {
    if (event.key === 'Enter') {
        handleUserLinkSubmit(userLink)
    }
  }

  return (
    <div className="App">
      <header className="App-header">
          <Container>

                  { topLinks.length > 0 &&
                      <ListGroup small>
                          <ListGroupItemHeading>
                              Top Links
                          </ListGroupItemHeading>

                          {topLinks.map((topLink) => {
                              return <ListGroupItem
                                    disabled={false}
                                    key={topLink.link}
                                    style={{textAlign: 'justify', fontSize: 14,}}
                                >
                                <Row>
                                    <Col style={{ textAlign: 'left' }}>{topLink.link}</Col>
                                    <Col style={{ textAlign: 'right' }}>{topLink.amount}</Col>
                                </Row>
                              </ListGroupItem>
                            })
                          }
                      </ListGroup>
                  }

                  <br/>

                  <InputGroup seamless onSubmit={handleUserLinkSubmit}>
                      <InputGroupAddon type="prepend">
                          <InputGroupText>
                              <FontAwesomeIcon icon={faLink} />
                          </InputGroupText>
                      </InputGroupAddon>


                      <FormInput
                          type={"url"}
                          valid={userLinkValid}
                          invalid={userLinkInvalid}
                          placeholder="https://..."
                          onChange={handleUserLinkChange}
                          onSubmit={handleUserLinkSubmit}
                          onKeyDown={handleUserLinkSubmitOnKey}
                      />

                      <InputGroupAddon type="append">
                          <Button onClick={() => handleUserLinkSubmit(userLink)}>
                              shrt!
                          </Button>
                      </InputGroupAddon>
                  </InputGroup>


                  <Collapse open={showLinkHash} style={{marginTop: '10px'}}>
                      <InputGroup>

                          <FormInput size="lg" value={urlToCopy} readOnly/>

                          <InputGroupAddon type="append" style={{align: 'center'}}>
                              <CopyToClipboard text={urlToCopy}>
                                <Button id="url-to-cpy-button" theme="info" onClick={() => setUrlToCopyCopied(true)}>cpy!</Button>
                              </CopyToClipboard>

                                <Popover
                                        placement="top"
                                        open={urlToCopyCopied}
                                        toggle={() => setUrlToCopyCopied(false)}
                                        target="#url-to-cpy-button"
                                        >
                                    <PopoverHeader>Cped!</PopoverHeader>
                                    <PopoverBody>
                                        Supposed to be in your clipboard..
                                    </PopoverBody>
                                </Popover>
                          </InputGroupAddon>
                      </InputGroup>
                  </Collapse>

                  <br/>
                  <br/>
                  <br/>
        <span style={{ marginTop: '50px' }}>Powered by <a href="https://github.com/lalabuy948/linkopus">linkopus</a></span>

        </Container>

      </header>
    </div>
  );
}

export default App;
