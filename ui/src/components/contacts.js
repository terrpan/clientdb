import * as React from 'react';
import {
  List,
  Datagrid,
  TextField,
  EmailField,
  ArrayField,
  SingleFieldList,
  ChipField,
  DateField,
  ReferenceField,
  Show,
  Edit,
  SimpleForm,  
  TextInput,
  ReferenceInput,
  SelectInput,
  ReferenceArrayInput,
  ArrayInput,
  SimpleFormIterator,
  useNotify, 
  useRefresh, 
  useRedirect,
  Create,
  TabbedForm,
  FormTab,
} from 'react-admin';


export const contactList = props => (
  <List {...props}>
    <Datagrid rowClick={"show"}>
      <TextField source="full_name" />
      <EmailField source="email" />
      <TextField source="phone_number" />
      <ArrayField source='client' label="Client">
        <SingleFieldList>
          <ChipField source="client_name" />
        </SingleFieldList>
      </ArrayField>
    </Datagrid>
  </List>
);

export const contactShow = props => (
  <Show {...props}>
    <TabbedForm>
      <FormTab label="Summary">
        <TextField source="id" />
        <TextField source="full_name" />
        <EmailField source="email" />
        <TextField source="phone_number" />
        <TextField source="role" />
        <DateField source="created_on" />
        <DateField source="modified_on" />
      </FormTab>
      <FormTab label="Client">
        <ArrayField source="client">
          <Datagrid>
            <ReferenceField source="id" reference="clients" link="show">
              <TextField source="client_name" />
            </ReferenceField>
          </Datagrid>
        </ArrayField>
      </FormTab>
    </TabbedForm>
  </Show>
);

export const contactEdit = props => (
  <Edit {...props}>
    <SimpleForm>
      <TextInput disabled source="id" />
      <TextInput source="first_name" />
      <TextInput source="last_name" />
      <TextInput disabled source="full_name" />
      <TextInput source="email" />
      <TextInput source="phone_number" />
      <TextInput source="role" />
      <ReferenceArrayInput source="attached_to_client" reference="clients" label="Client" allowEmpty>
        <ArrayInput>
        <SimpleFormIterator>
          <ReferenceInput source="client_id" reference="clients" label="Client" >
            <SelectInput optionText="client_name"/>
          </ReferenceInput>
        </SimpleFormIterator>
        </ArrayInput>
      </ReferenceArrayInput>
    </SimpleForm>
  </Edit>
);

export const ContactCreate = props => {
  const notify = useNotify();
  const redirect = useRedirect();
  const refresh = useRefresh();

  const onSuccess = () => {
    notify('Contact successfully created');
    redirect('/contacts');
    refresh();
  };

  return (
    <Create onSuccess={onSuccess} {...props}>
      <SimpleForm>
        <TextInput source="first_name" required/>
        <TextInput source="last_name" required/>
        <TextInput source="email" required/>
        <TextInput source="phone_number" />
        <TextInput source="role" />
        <ReferenceArrayInput source="attached_to_client" reference="clients" label="Client" allowEmpty>
          <ArrayInput>
          <SimpleFormIterator>
            <ReferenceInput source="client_id" reference="clients" label="Client" >
              <SelectInput optionText="client_name"/>
            </ReferenceInput>
          </SimpleFormIterator>
          </ArrayInput>
        </ReferenceArrayInput>
      </SimpleForm>
    </Create>
  );
}