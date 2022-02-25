import * as React from 'react';
import {
  List,
  Datagrid,
  TextField,
  ArrayField,
  DateField,
  Show,
  TabbedShowLayout,
  Tab,
  TextInput,
  Edit,
  SimpleForm,
  Create,
  UrlField,
  EmailField,
  useNotify, 
  useRefresh, 
  useRedirect,
  TabbedForm,
  FormTab,
  SelectInput,
  ReferenceInput, 
} from 'react-admin';

//TODO: not in use do to missing functionality in backend
const clientFilters = [
  <TextInput label="Search" source="q" alwaysOn />,
]

export const clientList = props => (
  <List {...props}>
    <Datagrid rowClick={"show"}>
      <TextField source="client_name" />
      <TextField source="slack_channel" />
      <UrlField source="web_url" />
    </Datagrid>
  </List>
);

export const clientShow = (props) => (
  <Show {...props}>
    <TabbedShowLayout>
    <Tab label="Summary">
      <TextField source="id" />
      <TextField source="client_name" />
      <TextField source="slack_channel" />
      <UrlField source="web_url" />
      <DateField source="modified_on" showTime/>
      <DateField source="created_on" showTime/>
    </Tab>
    <Tab label="Contacts">
      <ArrayField source="client_contacts">
        <Datagrid>
          <TextField source="full_name" />
          <EmailField source="email" />
          <TextField source="phone_number" />
          <TextField source="role" />
        </Datagrid>
      </ArrayField>
    </Tab>
    <Tab label="Services">
      <ArrayField source="managed_services">
        <Datagrid>
          <TextField source="service_name" />
          <TextField source="service_type" />
          <TextField source="service_description" />
        </Datagrid>
      </ArrayField>
    </Tab>
    </TabbedShowLayout>
  </Show>
);



export const clientEdit = props => (
  <Edit {...props}>
    <TabbedForm>
      <FormTab label="Summary">
        <TextInput disabled source="id" />
        <TextInput source='client_name'/>
        <TextInput source="slack_channel" />
        <TextInput source="web_url" />
      </FormTab>
      <FormTab label="Contacts">
        <ReferenceInput source="client_contacts" reference="contacts">
          {/* TODO: Create a custom post method to post an update to a special endpoint */}
          <SelectInput optionText="full_name"  />
        </ReferenceInput>
      </FormTab>
      <FormTab label="Services">
        <ReferenceInput source="client_services" reference="services">
          {/* TODO: Create a custom post method to post an update to a special endpoint */}
          <SelectInput optionText="service_name" />
        </ReferenceInput>
      </FormTab>
    </TabbedForm>
  </Edit>
);

export const ClientCreate = props => {
  const notify = useNotify();
  const redirect = useRedirect();
  const refresh = useRefresh();

  const onSuccess = () => {
    notify('Created successfully');
    redirect('/clients');
    refresh();
  };

  return (
    <Create onSuccess={onSuccess} {...props} >
      <SimpleForm>
      <TextInput source='client_name'/>
      <TextInput source="slack_channel" />
      <TextInput source="web_url" />
      </SimpleForm>
    </Create>
  );
}