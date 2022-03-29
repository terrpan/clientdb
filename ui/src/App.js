import * as React from 'react';
import { Admin, Resource,  ListGuesser, ShowGuesser, EditGuesser } from 'react-admin';
import jsonServerProvider from 'ra-data-json-server';
import clientIcon from '@material-ui/icons/Book';
import serviceIcon from '@material-ui/icons/SettingsApplications';
import contactsIcon from '@material-ui/icons/Contacts';

import { 
  clientList, 
  clientShow, 
  clientEdit, 
  ClientCreate 
} from './components/clients';
import { 
  serviceList, 
  serviceShow,
  serviceEdit,
} from './components/services';
import {
  contactList,
  contactShow,
  contactEdit,
  ContactCreate,
} from './components/contacts';

const dataProvider = jsonServerProvider('http://localhost:3000/api'); // /api will be proxied to http://localhost:8080
// const dataProvider = jsonServerProvider('http://localhost:3000/api');

// const dataProvider = jsonServerProvider('http://clientdb-api:8080/api');

// const dataProvider = jsonServerProvider('http://8b9b-89-253-120-215.ngrok.io/api');


const app = () => (
  <Admin dataProvider={dataProvider}>
    <Resource 
      name="clients" 
      list={clientList} 
      show={clientShow}
      edit={clientEdit}
      create={ClientCreate}
      icon={clientIcon} 
    />
    <Resource
      name="services" 
      list={serviceList} 
      show={serviceShow} 
      edit={serviceEdit}
      icon={serviceIcon}
    />
    <Resource 
      name="contacts" 
      list={contactList}
      show={contactShow}
      edit={contactEdit}
      create={ContactCreate}
      icon={contactsIcon} 
    />

  </Admin>
);

export default app;