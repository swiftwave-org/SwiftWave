<script setup>
import { reactive, ref } from 'vue';
import { toast } from 'vue-sonner'
import ModalDialog from '../components/ModalDialog.vue';
import { useMutation } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import FilledButton from '../components/FilledButton.vue';


const props = defineProps({
    callbackOnPop: {
        type: Function,
        required: false,
        default: () => { }
    }
});

const isModalOpen = ref(false);



const details = reactive({
    domainId: 0,
    mode: 'letsencrypt', // letsencrypt, custom
    sslFullChain: '',
    sslPrivateKey: ''
});

// Issue SSL
const {
    mutate: issueAutoSsl,
    onError: onAutoIssueSslError,
    loading: isAutoIssueSslLoading,
    onDone: onAutoIssueSslDone
} = useMutation(
    gql`
    mutation ($id: Uint!) {
      issueSSL(id: $id) {
        name
      }
    }
  `,
    {
        variables: {
            id: ''
        }
    }
);

onAutoIssueSslDone(() => {
    toast.success('SSL issue request submitted successfully');
    closeModal();
});

onAutoIssueSslError((err) => {
    toast.error(err.message);
});

// Add Custom SSL
const {
    mutate: addCustomSSL,
    loading: isAddCustomSSLLoading,
    onError: onAddCustomSSLError,
    onDone: onAddCustomSSLDone
} = useMutation(
    gql`
    mutation ($id: Uint!, $input: CustomSSLInput!) {
      addCustomSSL(id: $id, input: $input) {
        name
      }
    }
  `,
    {
        variables: {
            id: '',
            input: {
                fullChain: '',
                privateKey: ''
            }
        }
    }
);

onAddCustomSSLDone(() => {
    toast.success('Custom SSL added successfully. It can take few minutes to apply properly');
    closeModal();
});

onAddCustomSSLError((err) => {
    toast.error(err.message);
});

const issueSSL = () => {
    if (details.mode === 'letsencrypt') {
        issueAutoSsl({ id: details.domainId });
    } else {
        details.sslFullChain = details.sslFullChain.trim();
        details.sslPrivateKey = details.sslPrivateKey.trim();
        if (details.sslFullChain === '' || details.sslPrivateKey === '') {
            toast.error('Please fill all the required fields');
            return;
        }
        addCustomSSL({
            id: details.domainId, input: {
                fullChain: details.sslFullChain,
                privateKey: details.sslPrivateKey
            }
        });
    }
};

const openModal = (id) => {
    details.domainId = id;
    details.mode = 'letsencrypt';
    details.sslFullChain = '';
    details.sslPrivateKey = '';
    isModalOpen.value = true;
};
const closeModal = () => {
    isModalOpen.value = false;
    props.callbackOnPop();
};


defineExpose({
    openModal,
    closeModal
});
</script>
<template>
    <teleport to="body">
        <ModalDialog :close-modal="closeModal" :is-open="isModalOpen" width="lg">
            <template v-slot:header>Issue SSL</template>
            <template v-slot:body>
                <p>You can issue SSL via Let's Encrypt or add custom SSL certificate.</p>
                <form @submit.prevent="">
                    <!-- Issuing method -->
                    <div class="mt-4">
                        <label class="block text-sm font-medium text-gray-700" for="name">
                            SSL Issuing Method
                        </label>
                        <div class="mt-1">
                            <select v-model="details.mode"
                                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
                                <option value="letsencrypt">
                                    Auto Issue SSL via Let's Encrypt
                                </option>
                                <option value="custom">
                                    Custom SSL Certificate
                                </option>
                            </select>
                        </div>
                    </div>
                    <!-- Private Key -->
                    <div class="mt-4" v-if="details.mode === 'custom'">
                        <label class="block text-sm font-medium text-gray-700" for="ssh_private_key">
                            SSL Private Key (In OpenSSL Format)
                        </label>
                        <div class="mt-2">
                            <textarea id="ssh_private_key" v-model="details.sslPrivateKey"
                                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                                placeholder="-----BEGIN RSA PRIVATE KEY-----
                                
....

-----END RSA PRIVATE KEY-----
                                " type="text" rows="5" />
                        </div>
                    </div>
                    <!-- certificate full chain -->
                    <div class="mt-4" v-if="details.mode === 'custom'">
                        <label class="block text-sm font-medium text-gray-700" for="ssh_private_key">
                            SSL Certificate FullChain (In OpenSSL Format)
                        </label>
                        <div class="mt-2">
                            <textarea id="ssh_private_key" v-model="details.sslFullChain"
                                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                                placeholder="-----BEGIN CERTIFICATE-----

....

-----END CERTIFICATE-----

-----BEGIN CERTIFICATE-----

.....

-----END CERTIFICATE-----
" type="text" rows="5" />
                        </div>
                    </div>
                </form>
            </template>
            <template v-slot:footer>
                <FilledButton class="w-full" type="primary" :loading="isAutoIssueSslLoading || isAddCustomSSLLoading"
                    :click="issueSSL">{{
                        details.mode === 'letsencrypt' ? ' Verify DNS And Issue SSL Certificate' : 'Add Custom SSL & Apply'
                    }}</FilledButton>
            </template>
        </ModalDialog>
    </teleport>
</template>