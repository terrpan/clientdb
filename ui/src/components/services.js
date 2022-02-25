import * as React from 'react';
import { 
  List, 
  Datagrid, 
  TextField,
  DateField, 
  Show,
  ArrayField,
  SingleFieldList,
  ChipField,
  ReferenceField,
  Edit,
  SimpleForm,
  TextInput,
  ReferenceInput,
  SelectInput,
  ReferenceArrayInput,
  ArrayInput,
  SimpleFormIterator,
  NumberInput,
  TabbedForm,
  FormTab,
} from 'react-admin';


export const serviceList = props => (
  <List {...props}>
    <Datagrid rowClick={"show"}>
      <TextField source="service_name" />
      <TextField source="service_type" />
      <ArrayField source='client' label="Clients">
        <SingleFieldList>
          <ChipField source="client_name" />
        </SingleFieldList>
      </ArrayField>
    </Datagrid>
  </List>
);

export const serviceShow = props => (
  <Show {...props}>
    <TabbedForm>
      <FormTab label="Summary">
        <TextField source="id" />
        <TextField source="service_name" />
        <TextField source="service_description" />
        <TextField source="service_type" />
        <DateField source="created_on" />
        <DateField source="modified_on" />
      </FormTab>
      <FormTab label="Clients">
      <ArrayField source="client">
          <Datagrid>
            <ReferenceField source="id" reference="clients" link="show">
            <TextField source="client_name" />
            </ReferenceField>
          </Datagrid>
        </ArrayField>
      </FormTab>
      <FormTab label="Economics">
        <TextField source="invoice_amount" />
        <TextField source="invoice_frequency" />
        <TextField source="management_fee" />
      </FormTab>
    </TabbedForm>
  </Show>
);

export const serviceEdit = props => (
  <Edit {...props}>
    <SimpleForm>
      <TextInput source="service_name" />
      <TextInput source="service_type" />
      <TextInput source="service_owner" />
      <TextInput source="service_description" />
      <TextInput source="service_status" />
      <TextInput source="invoice_frequency" />
      <ReferenceArrayInput source="attached_to_client" reference="clients" label="Client" allowEmpty>
        <ArrayInput>
        <SimpleFormIterator>
          <ReferenceInput source="client_id" reference="clients" label="Client" >
            <SelectInput optionText="client_name"/>
          </ReferenceInput>
        </SimpleFormIterator>
        </ArrayInput>
      </ReferenceArrayInput>
      <NumberInput source="invoice_amount" />
      <NumberInput source="management_fee" />
    </SimpleForm>
  </Edit>
);