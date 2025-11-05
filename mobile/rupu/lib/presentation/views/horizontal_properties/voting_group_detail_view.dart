import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:rupu/domain/entities/voting_results_result.dart';
import 'package:rupu/presentation/views/horizontal_properties/controllers/voting_group_detail_controller.dart';

class VotingGroupDetailView extends StatelessWidget {
  final int propertyId;
  final int votingGroupId;

  const VotingGroupDetailView({
    super.key,
    required this.propertyId,
    required this.votingGroupId,
  });

  @override
  Widget build(BuildContext context) {
    final controller = Get.find<VotingGroupDetailController>();
    return Scaffold(
      appBar: AppBar(
        title: Obx(() => Text(controller.group.value?.name ?? 'Detalle de grupo')),
      ),
      body: Obx(
        () {
          if (controller.isLoading.value) {
            return const Center(child: CircularProgressIndicator());
          }
          if (controller.errorMessage.value != null) {
            return Center(child: Text(controller.errorMessage.value!));
          }
          return ListView(
            children: [
              Text('Resultados de la votación'),
              // TODO: Display voting results
            ],
          );
        },
      ),
    );
  }
}
