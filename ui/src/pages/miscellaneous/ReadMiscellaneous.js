import React, { useEffect, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';

import { Grid, Button } from '@material-ui/core';

import PageTitle from '../../components/PageTitle/PageTitle';
import Dialog from '../../components/Dialog/Dialog';

import CustomTable from '../../components/Table/Table.js';

import { formattedDate } from '../../utils/helpers';

import CreateUpdateMiscellaneous from './CreateUpdateMiscellaneous';
import {
  miscellaneousRead,
  miscellaneousCreate,
  miscellaneousUpdate,
  miscellaneousDelete,
} from '../../store/actions/miscellaneous';

const headCells = [
  {
    id: 'id',
    numeric: false,
    disablePadding: false,
    label: 'SN',
  },
  {
    id: 'description',
    numeric: false,
    disablePadding: false,
    label: 'Description',
  },
  {
    id: 'price',
    numeric: false,
    disablePadding: false,
    label: 'Price',
  },
  {
    id: 'date',
    numeric: false,
    disablePadding: false,
    label: 'Date',
  },
];

const ReadMiscellaneous = () => {
  let rows = [];
  const dispatch = useDispatch();

  const [openCreateUpdateMiscellaneous, setOpenCreateUpdateMiscellaneous] =
    useState(false);
  const [valueForm, setValueForm] = useState(null);

  const miscellaneousItems = useSelector(
    (state) => state.miscellaneous.miscellaneous,
  );
  const pageNo = useSelector((state) => state.miscellaneous.pageNo);
  const rowsPerPage = useSelector((state) => state.miscellaneous.rowsPerPage);
  const orderBy = useSelector((state) => state.miscellaneous.orderBy);
  const order = useSelector((state) => state.miscellaneous.order);
  const token = useSelector((state) => state.user.token);
  const isLoading = useSelector((state) => state.miscellaneous.formLoading);
  const totalCount = useSelector((state) => state.miscellaneous.count);

  if (miscellaneousItems) {
    rows = miscellaneousItems.map((miscellaneous, idx) => {
      return {
        id: miscellaneous.id ? miscellaneous.id : idx + 1,
        description: miscellaneous.description
          ? miscellaneous.description
          : 'Not Specified',
        price: miscellaneous.price ? miscellaneous.price : 'Not Specified',
        date: miscellaneous.date
          ? formattedDate(miscellaneous.date)
          : 'Not Specified',
      };
    });
  }

  useEffect(() => {
    document.title = `Miscellaneous | ${process.env.REACT_APP_NAME}`;
  }, []);

  const handleRequestSort = (property) => {
    const isAsc = orderBy === property && order === 'asc';
    const direction = isAsc ? 'desc' : 'asc';
    dispatch(
      miscellaneousRead({
        token,
        pageNo,
        rowsPerPage,
        order: direction,
        orderBy: property,
      }),
    );
  };

  const handleChangePage = (_event, newPage) => {
    dispatch(
      miscellaneousRead({
        token,
        pageNo: newPage,
        rowsPerPage,
        orderBy,
        order,
      }),
    );
  };

  useEffect(() => {
    if (miscellaneousItems.length === 0 && pageNo >= 0) {
      handleChangePage.call(null, +pageNo);
    }
  }, [pageNo, miscellaneousItems]);

  const handleChangeRowsPerPage = (event) => {
    dispatch(
      miscellaneousRead({
        token,
        pageNo,
        rowsPerPage: parseInt(event.target.value, 10),
        orderBy,
        order,
      }),
    );
  };

  const handleOpenCreateMiscellaneous = () => {
    setOpenCreateUpdateMiscellaneous(true);
  };

  const handleOpenEditMiscellaneous = (values) => {
    setValueForm(values);
    setOpenCreateUpdateMiscellaneous(true);
  };

  const handleCloseCreateOrEditMiscellaneous = () => {
    setValueForm(null);
    setOpenCreateUpdateMiscellaneous(false);
  };

  const handleSubmitCreateUpdateMiscellaneous = (formValues) => {
    formValues.date = new Date(formValues.date).toISOString();
    if (formValues.id) {
      dispatch(
        miscellaneousUpdate({
          formValues,
          token,
        }),
      );
    } else {
      dispatch(
        miscellaneousCreate({
          formValues,
          token,
        }),
      );
    }
    handleCloseCreateOrEditMiscellaneous();
  };

  const handleSubmitDeleteMiscellaneous = (id) => {
    const result = window.confirm(
      `Are you sure you want to delete this miscellaneous item?`,
    );
    if (result) {
      dispatch(
        miscellaneousDelete({
          id,
          token,
          // pageNo,
          // rowsPerPage,
          // orderBy,
          // order,
          // totalCount,
        }),
      );
    }
  };

  return (
    <>
      <PageTitle
        title={'Miscellaneous'}
        button={
          <Button
            variant='outlined'
            size='medium'
            color='secondary'
            onClick={handleOpenCreateMiscellaneous}
          >
            Add Miscellaneous
          </Button>
        }
      />
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <CustomTable
            tableTitle='All Miscellaneous'
            order={order}
            orderBy={orderBy}
            requestSort={handleRequestSort}
            headCells={headCells}
            rows={rows}
            openEditMiscellaneous={handleOpenEditMiscellaneous}
            submitDeleteMiscellaneous={handleSubmitDeleteMiscellaneous}
            totalCount={totalCount}
            pageNo={pageNo}
            rowsPerPage={rowsPerPage}
            changePage={handleChangePage}
            changeRowsPerPage={handleChangeRowsPerPage}
          />
        </Grid>
      </Grid>
      <Dialog
        fullWidth={true}
        maxWidth='xs'
        open={openCreateUpdateMiscellaneous}
        handleClose={handleCloseCreateOrEditMiscellaneous}
        title={`${valueForm ? 'Edit' : 'Create'} Miscellaneous`}
        handleSubmit={handleSubmitCreateUpdateMiscellaneous}
        initialValues={valueForm}
        isLoading={isLoading}
      >
        <CreateUpdateMiscellaneous initialValues={valueForm} />
      </Dialog>
    </>
  );
};

export default ReadMiscellaneous;
